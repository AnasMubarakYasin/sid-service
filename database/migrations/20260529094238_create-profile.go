package migrations

import (
	"context"
	"fmt"
	"sid/service/model"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		// Create a single table
		_, err := db.NewCreateTable().
			Model((*model.Profile)(nil)).
			IfNotExists().
			Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		// Drop table
		_, err := db.NewDropTable().
			Model((*model.Profile)(nil)).
			IfExists().
			Cascade(). // Drop dependent objects
			Exec(ctx)
		return err
	})
}
