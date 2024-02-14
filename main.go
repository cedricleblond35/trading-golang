package main

import (
	"fmt"
	"log"
	"time"

	"trading/cmd/candleCollector"
	"trading/internal/config"

	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{
		Use:     "main",
		Version: config.Version(),
	}
	c.AddCommand(candleCollector.Command)
	c.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Version for trading",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(c.Version)
		},
	})

	if err := c.Execute(); err != nil {
		defer func() {
			time.Sleep(5 * time.Second)
		}()

		log.Fatalf("%+v", err)
	}
}
