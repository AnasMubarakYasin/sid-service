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
		_, err := db.NewCreateTable().
			Model((*model.Experiment)(nil)).
			IfNotExists().
			Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			Model((*model.Experiment)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		return err
	})
}
