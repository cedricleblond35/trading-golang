package xtb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"trading/internal/database"
	"trading/internal/model"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type XTB struct {
	URL    url.URL
	Conn   *websocket.Conn
	status bool
	PDB    *database.GORM
}

func NewXTB(pdb *database.GORM) *XTB {
	return &XTB{
		PDB: pdb,
	}
}

type Result struct {
	Candles []byte
	Period  int
}

func (xtb *XTB) Connection(host, path string) error {
	if host == "" {
		fmt.Println("error host")
	}
	if path == "" {
		fmt.Println("error path")
	}
	u := url.URL{
		Scheme: "wss",
		Host:   host,
		Path:   path,
	}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		return err
	}
	xtb.Conn = c

	return nil
}

func (xtb *XTB) Login(user, password string) error {
	jsonData, err := json.Marshal(Login{
		Command: "login",
		Arguments: ArgsLogin{
			UserID:   user,
			Password: password,
		},
	})
	if err != nil {
		return err
	}

	if err := xtb.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return err
	}
	_, message, err := xtb.Conn.ReadMessage()
	if err != nil {
		return err
	}

	var rc ResponseConn
	err = json.Unmarshal([]byte(message), &rc)
	if err != nil {
		return err
	}
	xtb.status = rc.Status
	fmt.Println("connected to xtb:", rc.Status)
	return nil
}

func takeCandle(messageOut chan Result, interrupt chan struct{}, ticker *time.Ticker, xtb *XTB, period, offset int) {
	defer close(interrupt)
	for t := range ticker.C {
		fmt.Println("tick à :", t, " period", period)
		if xtb.status {
			timestamp := time.Now().UTC().UnixMilli()
			fmt.Println("demande de bougie : ", timestamp)

			jsonDataChartRangeRequest, err := json.Marshal(GetChartLastRequest{
				Command: "getChartLastRequest",
				Arguments: ArgsGetChartLastRequest{
					ArgsInfogetChartLastRequest{
						Period: period,
						Start:  timestamp - int64(offset)*1000,
						Symbol: "US100",
					},
				},
			})
			if err != nil {
				fmt.Println("erreur Marshal:", err)
				// quit <- os.Kill
			}

			if err := xtb.Conn.WriteMessage(websocket.TextMessage, jsonDataChartRangeRequest); err != nil {
				fmt.Println("erreur Marshal:", err)
				log.Println(err)
			}
			messageType, message, err := xtb.Conn.ReadMessage()
			if err != nil {
				fmt.Println("erreur ReadMessage:", err)
				xtb.Conn.Close()
				break
			}

			if messageType == websocket.TextMessage {
				res := new(Result)
				res.Candles = message
				res.Period = period
				messageOut <- *res
			} else if messageType == websocket.BinaryMessage {
				// handle binary message
			} else if messageType == websocket.CloseMessage {
				// handle close message
			}
		}
	}
}

func (xtb *XTB) Collected(symbol string, period int) error {
	ticker := time.NewTicker(15 * time.Second)
	messageOut := make(chan Result, 1)
	interrupt := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	offset := 60 * 5 // 5m
	go takeCandle(messageOut, interrupt, ticker, xtb, 1, offset)
	offset = 60 * 30 // 30 min
	go takeCandle(messageOut, interrupt, ticker, xtb, 5, offset)
	offset = 60 * 45 // 45 min
	go takeCandle(messageOut, interrupt, ticker, xtb, 15, offset)
	offset = 60 * 60 * 12 // 12h
	go takeCandle(messageOut, interrupt, ticker, xtb, 240, offset)

	for {
		select {
		case message := <-messageOut:
			var respData ResponseChartLastRequest
			err := json.Unmarshal(message.Candles, &respData)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			err = insertData(respData, message.Period, xtb)
			if err != nil {
				return err
			}

		case <-quit:
			fmt.Println("Received interruption signal, gracefully shutting down")
			err := xtb.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {

				log.Println("write close:", err)

				return nil

			}

			return nil
		case <-interrupt:

			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then

			// waiting (with timeout) for the server to close the connection.

			err := xtb.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {

				log.Println("write close:", err)

				return nil

			}

			return nil

		}
	}
}

func insertData(resp ResponseChartLastRequest, period int, xtb *XTB) error {
	// resp.ReturnData.RateInfos = resp.ReturnData.RateInfos[:len(resp.ReturnData.RateInfos)-1]
	for _, ligne := range resp.ReturnData.RateInfos {
		x := model.NewCandle()
		if err := xtb.PDB.LoadLast(x, "period = ? AND ctm =?", period, ligne.Ctm/1000); errors.Is(err, gorm.ErrRecordNotFound) {
			open := ligne.Open / 100
			high := (ligne.Open + ligne.High) / 100
			low := (ligne.Open + ligne.Low) / 100
			close := (ligne.Open + ligne.Close) / 100

			m := model.NewCandle()
			m.Open = sql.NullFloat64{Valid: true, Float64: open}
			m.Close = sql.NullFloat64{Valid: true, Float64: close}
			m.Low = sql.NullFloat64{Valid: true, Float64: low}
			m.High = sql.NullFloat64{Valid: true, Float64: high}
			m.Period = sql.NullInt16{Valid: true, Int16: int16(period)}
			m.Ctm = sql.NullInt64{Valid: true, Int64: ligne.Ctm / 1000}
			m.Date = ligne.CtmString
			m.Vol = sql.NullInt32{Valid: true, Int32: int32(ligne.Vol)}
			err := xtb.PDB.Save(m)
			if err != nil {
				return errors.Wrap(err, "could not create candle record")
			}
		}
	}
	return nil
}
