package migrations

import (
	"context"
	"fmt"
	"sid/service/feature"
	"sid/service/model"
	"sid/service/repository"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		// Create a single table
		_, err := db.NewCreateTable().
			Model((*model.User)(nil)).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}
		users := []*feature.ParamSignUp{
			{Username: "teacher", Role: "teacher", Password: "teacher", Confirm: "teacher"},
			{Username: "roger", Role: "teacher", Password: "roger", Confirm: "roger"},
			{Username: "student", Role: "student", Password: "student", Confirm: "student"},
			{Username: "alice", Role: "student", Password: "alice", Confirm: "alice"},
		}
		ru := repository.NewUser(db)
		rp := repository.NewProfile(db)
		f := feature.NewUser(ru, rp)
		for _, v := range users {
			_, err := f.UserSignUp(ctx, v)
			if err != nil {
				return err
			}
		}
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		// Drop table
		_, err := db.NewDropTable().
			Model((*model.User)(nil)).
			IfExists().
			Cascade(). // Drop dependent objects
			Exec(ctx)
		return err
	})
}
