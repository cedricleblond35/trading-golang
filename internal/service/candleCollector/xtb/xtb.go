package xtb

import (
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
	fmt.Println("connecté:", rc.Status)
	return nil
}

func (xtb *XTB) Collected(symbol string, period int) error {
	ticker := time.NewTicker(1 * time.Second)
	messageOut := make(chan []byte, 1)
	interrupt := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		defer close(interrupt)
		for t := range ticker.C {
			fmt.Println("tick à :", t)
			fmt.Println("status :", xtb.status)
			fmt.Println("PERIOD :", period)
			if xtb.status {
				timestamp := time.Now().UTC().UnixMilli()
				fmt.Println("demande de bougie : ", timestamp)

				jsonDataChartRangeRequest, err := json.Marshal(GetChartLastRequest{
					Command: "getChartLastRequest",
					Arguments: ArgsGetChartLastRequest{
						ArgsInfogetChartLastRequest{
							Period: 1440,
							Start:  timestamp - (60*60*24*30)*1000,
							Symbol: "US100",
						},
					},
				})
				// jsonDataChartRangeRequest, err := json.Marshal(GetChartRangeRequest{
				// 	Command: "getChartRangeRequest",
				// 	Arguments: ArgsGetChartRangeRequest{
				// 		ArgsInfoGetChartRangeRequest{
				// 			Start:  timestamp - (60 * 60 * 1000),
				// 			End:    timestamp,
				// 			Period: 1,
				// 			Symbol: "US100",
				// 			Ticks:  0,
				// 		},
				// 	},
				// })
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
					messageOut <- message
				} else if messageType == websocket.BinaryMessage {
					// handle binary message
				} else if messageType == websocket.CloseMessage {
					// handle close message
				}
			}
		}
	}()

	m := model.NewCandle()
	for {
		select {
		case message := <-messageOut:
			var respData ResponseChartLastRequest
			err := json.Unmarshal([]byte(message), &respData)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			info := respData.ReturnData.RateInfos
			for i, ligne := range info {
				open := ligne.Open
				high := ligne.Open + ligne.High
				low := ligne.Open + ligne.Low
				close := ligne.Open + ligne.Close
				s := fmt.Sprintf("%.0f", ligne.Open)
				fmt.Println("s:", s)
				pivotCamarilla := model.NewPivot()
				pivotCamarilla.PivotCamarilla(high, low, close)
				fmt.Println(i, "date:", ligne.CtmString, "open", open, " | close", ligne.Close)
				fmt.Println("pivotCamarilla R1:", pivotCamarilla.R1, " R2:", pivotCamarilla.R2, " R3:", pivotCamarilla.R3)

				err := xtb.PDB.Create(m)
				m.Open = float32(open)
				m.Close = float32(close)
				m.Low = float32(low)
				m.High = float32(high)
				m.Period = 1440
				// m.Ctm = ligne.Ctm

				if err != nil {
					return errors.Wrap(err, "could not create candle record")
				}
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
