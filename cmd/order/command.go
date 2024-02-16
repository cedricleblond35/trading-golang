package order

import (
	"context"
	"os"

	"trading/internal/database"
	"trading/internal/xtb"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:          "order",
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

	return nil
}
