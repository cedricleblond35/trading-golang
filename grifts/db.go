package grifts

import (
	"context"
	"fmt"

	"trading/internal/database"

	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {
	_ = grift.Desc("reset", "Reset the databases")
	_ = grift.Add("reset", func(c *grift.Context) error {
		ctx := context.Background()
		pdb, err := database.GormOpen(ctx, true)
		if err != nil {
			return errors.Wrap(err, "could not connect to database")
		}

		tables := []string{"user"}

		for _, tablename := range tables {
			err = pdb.DB().Exec(fmt.Sprintf("DELETE FROM %s", tablename)).Error
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})

	_ = grift.Desc("seed", "Seeds the database with basic data")
	_ = grift.Add("seed", func(c *grift.Context) error {
		ctx := context.Background()
		pdb, err := database.GormOpen(ctx, true)
		if err != nil {
			return errors.Wrap(err, "could not connect to database")
		}

		fmt.Println("Creating user..")
		user := CreateFakeUser()
		err = pdb.Create(user)
		if err != nil {
			return errors.Wrap(err, "could not create user record")
		}

		defer fmt.Printf("=======> user#%d created with name: %s\n", user.ID, user.Name)

		return nil
	})
})
