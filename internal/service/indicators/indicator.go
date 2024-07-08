package indicator

import (
	"context"
	"fmt"
	"time"

	"trading/internal/database"
	"trading/internal/model"

	"github.com/redis/go-redis/v9"
)

type Indicator struct {
	ctx     context.Context
	pdb     *database.GORM
	redis   *redis.Client
	candles []model.CandleUS100
	candle  model.CandleUS100
}

func NewIndicator(ctx context.Context, pdb *database.GORM, rdb *redis.Client) *Indicator {
	return &Indicator{
		ctx:   ctx,
		pdb:   pdb,
		redis: rdb,
	}
}

func (indicator *Indicator) PivotCamarilla(high, low, close float64) error {
	pc := model.NewPivot()
	pc.PivotCamarilla(high, low, close)
	err := indicator.redis.Set(indicator.ctx, "pc.r1", pc.R1, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pc.r2", pc.R2, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pc.r3", pc.R3, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pc.r4", pc.R4, 0).Err()
	if err != nil {
		panic(err)
	}

	err = indicator.redis.Set(indicator.ctx, "pc.s1", pc.S1, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pc.s2", pc.S2, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pc.s3", pc.S3, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pc.s4", pc.S4, 0).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (indicator *Indicator) PivotWoodie(high, low, close float64) error {
	pw := model.NewPivot()
	pw.PivotWoodie(high, low, close)
	err := indicator.redis.Set(indicator.ctx, "pw.p", pw.P, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pw.r1", pw.R1, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pw.r2", pw.R2, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pw.s1", pw.S1, 0).Err()
	if err != nil {
		panic(err)
	}
	err = indicator.redis.Set(indicator.ctx, "pw.s2", pw.S2, 0).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (indicator *Indicator) Calcul() error {
	timestamp := time.Now().Add(time.Hour*time.Duration(-24)).UTC().UnixMilli() / 1000
	fmt.Println("time:", timestamp)
	err := indicator.pdb.LoadLast(&indicator.candle, "period = ? AND ctm < ?", 1440, timestamp)
	if err != nil {
		return err
	}

	fmt.Println("============>date:", indicator.candle)

	// err = indicator.PivotCamarilla(indicator.candle.High.Float64, indicator.candle.Low.Float64, indicator.candle.Close.Float64)
	// if err != nil {
	// 	return err
	// }

	// r1 := indicator.redis.Get(indicator.ctx, "pc.r1")
	// fmt.Println("Camarilla R1:", r1)
	// err = indicator.PivotWoodie(indicator.candle.High.Float64, indicator.candle.Low.Float64, indicator.candle.Close.Float64)
	// if err != nil {
	// 	return err
	// }

	err = indicator.pdb.LoadsDESC(&indicator.candles, 100, -1, "period = ?", 5)
	if err != nil {
		return err
	}

	// fmt.Printf("indicator.candles: %+v\n", indicator.candles)
	// fmt.Println("--------------------------------------------------------------------")
	candles := reverse(indicator.candles)
	fmt.Printf("candles inver: %+v\n", candles)

	st := model.NewSupertrend()
	st.Calcul(candles)
	fmt.Println("---------------->", st.Get())

	return nil
}

func reverse(numbers []model.CandleUS100) []model.CandleUS100 {
	newCandles := make([]model.CandleUS100, 0, len(numbers))
	for i := len(numbers) - 1; i >= 0; i-- {
		newCandles = append(newCandles, numbers[i])
	}
	return newCandles
}
