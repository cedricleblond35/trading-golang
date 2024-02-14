package candleCollector

import (
	"context"
	"fmt"

	"trading/internal/database"
	candlecollectorPkg "trading/internal/service/candleCollector"

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

	pdb, err := database.GormOpen(ctx, false)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	err = candlecollectorPkg.Process(ctx, pdb)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
