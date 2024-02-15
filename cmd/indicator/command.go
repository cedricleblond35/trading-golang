package indicator

import (
	"context"
	"fmt"

	"trading/internal/database"
	"trading/internal/redis"

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

	_, err := database.GormOpen(ctx, false)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	rdb, err := redis.NewRedis()
	if err != nil {
		return errors.Wrap(err, "could not create redis client")
	}

	// test
	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	return nil
}
