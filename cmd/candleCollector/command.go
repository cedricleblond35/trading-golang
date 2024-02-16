package candleCollector

import (
	"context"
	"os"

	"trading/internal/database"
	"trading/internal/redis"
	candlecollectorPkg "trading/internal/service/candleCollector"
	"trading/internal/xtb"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:          "candleCollector",
	Short:        "",
	Args:         cobra.NoArgs,
	RunE:         action,
	SilenceUsage: true,
}

func action(c *cobra.Command, _ []string) error {
	ctx := context.Background()

	// Database
	pdb, err := database.GormOpen(ctx, false)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	// redis
	rdb, err := redis.NewRedis()
	if err != nil {
		return errors.Wrap(err, "could not create redis client")
	}

	// Connection XTB
	xtbConn := xtb.NewXTB()
	err = xtbConn.Connection("ws.xtb.com", "/real")
	if err != nil {
		return err
	}

	err = xtbConn.Login(os.Getenv("XTB_USER"), os.Getenv("XTB_PWD"))
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, "xtbConnStatus", xtbConn.Status, 0).Err()
	if err != nil {
		panic(err)
	}

	// err = rdb.Set(ctx, "xtbConn", xtbConn.Conn, 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// candle
	cc := candlecollectorPkg.NewCandleCollector(ctx, pdb, xtbConn)
	err = cc.Collected("US100")
	if err != nil {
		return err
	}

	return nil
}
