package candlecollector

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"trading/internal/database"
	"trading/internal/model"
	"trading/internal/xtb"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Result struct {
	Candles []byte
	Period  int
}

type CandleCollector struct {
	ctx     context.Context
	pdb     *database.GORM
	xtbConn *xtb.XTB
	redis   *redis.Client
}

func NewCandleCollector(ctx context.Context, pdb *database.GORM, xtbConn *xtb.XTB, rdb *redis.Client) *CandleCollector {
	return &CandleCollector{
		ctx:     ctx,
		pdb:     pdb,
		xtbConn: xtbConn,
		redis:   rdb,
	}
}

func takeCandle(messageOut chan Result, interrupt chan struct{}, ticker *time.Ticker, xtb *xtb.XTB, period, offset int) {
	defer close(interrupt)
	for t := range ticker.C {
		fmt.Println("tick à :", t, " period", period)
		if xtb.Status {
			timestamp := time.Now().UTC().UnixMilli()
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
			} else if messageType == websocket.CloseMessage {
				fmt.Println("message de contrôle de fermeture:", messageType) // TODO: Gérer le cas
			}
		}
	}
}

func (cc *CandleCollector) Collected(symbol string) error {
	ticker := time.NewTicker(15 * time.Second)
	messageOut := make(chan Result, 1)
	interrupt := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	offset := 60 * 3 // 5m
	go takeCandle(messageOut, interrupt, ticker, cc.xtbConn, 1, offset)
	offset = 60 * 15 // 30 min
	go takeCandle(messageOut, interrupt, ticker, cc.xtbConn, 5, offset)
	offset = 60 * 45 // 45 min
	go takeCandle(messageOut, interrupt, ticker, cc.xtbConn, 15, offset)
	offset = 60 * 60 * 12 // 12h
	go takeCandle(messageOut, interrupt, ticker, cc.xtbConn, 240, offset)
	offset = 60 * 60 * 24 * 60 // 12h
	go takeCandle(messageOut, interrupt, ticker, cc.xtbConn, 1440, offset)

	for {
		select {
		case message := <-messageOut:
			var respData ResponseChartLastRequest
			err := json.Unmarshal(message.Candles, &respData)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			err = insertCandles(respData, message.Period, cc.pdb)
			if err != nil {
				return err
			}

		case <-quit:
			fmt.Println("Received interruption signal, gracefully shutting down")
			err := cc.xtbConn.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {

				log.Println("write close:", err)

				return nil

			}

			return nil
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := cc.xtbConn.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return nil
			}

			return nil
		}
		val, err := cc.redis.Get(cc.ctx, "pc.r1").Result()
		if err != nil {
			return err
		}
		fmt.Println("*************pc.r1", val)
		val, err = cc.redis.Get(cc.ctx, "xxxxc.r1").Result()
		if err != nil {
			fmt.Println("demande de calcul")
		}
		fmt.Println("*************xxxxc.r1", val)
		val, err = cc.redis.Get(cc.ctx, "pw.r2").Result()
		if err != nil {
			return err
		}
		fmt.Println("************* pw.r2", val)
	}
}

func insertCandles(resp ResponseChartLastRequest, period int, db *database.GORM) error {
	// resp.ReturnData.RateInfos = resp.ReturnData.RateInfos[:len(resp.ReturnData.RateInfos)-1]
	for _, ligne := range resp.ReturnData.RateInfos {
		open := ligne.Open / 100
		high := (ligne.Open + ligne.High) / 100
		low := (ligne.Open + ligne.Low) / 100
		cl := (ligne.Open + ligne.Close) / 100

		m := model.NewCandle()
		if err := db.LoadLast(m, "period = ? AND ctm =?", period, ligne.Ctm/1000); errors.Is(err, gorm.ErrRecordNotFound) {

			// m := model.NewCandle()
			m.Open = sql.NullFloat64{Valid: true, Float64: open}
			m.Close = sql.NullFloat64{Valid: true, Float64: cl}
			m.Low = sql.NullFloat64{Valid: true, Float64: low}
			m.High = sql.NullFloat64{Valid: true, Float64: high}
			m.Ctm = sql.NullInt64{Valid: true, Int64: ligne.Ctm / 1000}
			m.Date = ligne.CtmString
			m.Vol = sql.NullInt32{Valid: true, Int32: int32(ligne.Vol)}
			m.Period = sql.NullInt16{Valid: true, Int16: int16(period)}
			err := db.Save(m)
			if err != nil {
				return errors.Wrap(err, "could not create candle record")
			}
		} else {
			values := map[string]interface{}{
				"close": sql.NullFloat64{Valid: true, Float64: cl},
				"high":  sql.NullFloat64{Valid: true, Float64: high},
				"low":   sql.NullFloat64{Valid: true, Float64: low},
				"vol":   sql.NullInt32{Valid: true, Int32: int32(ligne.Vol)},
			}
			err := db.Update(m, values)
			if err != nil {
				return errors.Wrap(err, "could not update candle record")
			}
		}
	}
	return nil
}
