package indicator

import (
	"context"

	"trading/internal/database"
	"trading/internal/redis"
	indicatorsPkg "trading/internal/service/indicators"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:          "indicator",
	Short:        "calculation of indicators, example SMA, Pivot ...",
	Args:         cobra.NoArgs,
	RunE:         action,
	SilenceUsage: true,
}

func action(c *cobra.Command, _ []string) error {
	ctx := context.Background()

	pdb, err := database.GormOpen(ctx, false)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	rdb, err := redis.NewRedis()
	if err != nil {
		return errors.Wrap(err, "could not create redis client")
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////
	indicators := indicatorsPkg.NewIndicator(ctx, pdb, rdb)
	err = indicators.Calcul()
	if err != nil {
		return err
	}

	return nil
}
