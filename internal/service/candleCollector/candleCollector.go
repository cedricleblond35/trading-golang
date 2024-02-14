package candlecollector

import (
	"context"
	"os"

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

	err = x.Login(os.Getenv("USER"), os.Getenv("PWD"))
	if err != nil {
		return err
	}

	err = x.Collected("US100", int(model.PERIOD_M1))
	if err != nil {
		return err
	}

	return nil
}
