package candlecollector

import (
	"context"

	"trading/internal/database"
	"trading/internal/model"
	"trading/internal/service/candleCollector/xtb"
)

type Controller struct{}

func Process(ctx context.Context, pdb *database.GORM) error {
	x := xtb.NewXTB(pdb)
	err := x.Connection("ws.xtb.com", "/real")
	if err != nil {
		return err
	}

	err = x.Login("1502064", "1976Drick!")
	if err != nil {
		return err
	}

	err = x.Collected("US100",int(model.PERIOD_M5))
	if err != nil {
		return err
	}



	return nil
}
