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
		_, err := db.NewCreateTable().
			Model((*model.Subject)(nil)).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}
		params := []*feature.ParamCreateSubject{
			{Title: "Fenomena Listrik Dinamis", Description: "Listrik Dinamis ↓ Muatan listrik bergerak ↓ Arus listrik ↓ Elektron bebas ↓ Beda potensial ↓ Rangkaian listrik sederhana", Video: "https://www.youtube.com/watch?v=VrpdtzA8MhI&t=13s", Questions: []string{"Mengapa arus listrik dapat mengalir?", "Apa fungsi baterai?", "Mengapa lampu padam ketika rangkaian terbuka?"}},
			{Title: "Hukum Ohm", Description: "Tegangan Listrik ↓ Arus Listrik ↓ Hambatan Listrik ↓ Hukum Ohm ↓ Hubungan V, I, dan R ↓ Penerapan dalam Kehidupan", Video: "https://www.youtube.com/watch?v=aptOXSNpDp4&t=45s", Questions: []string{"Mengapa lampu bisa lebih terang?", "Mengapa arus listrik bisa berubah?", "Apa yang menyebabkan alat elektronik bekerja lebih kuat?", "Bagaimana hubungan tegangan dan arus?"}},
		}
		rs := repository.NewSubject(db)
		f := feature.NewSubject(rs)
		for _, v := range params {
			_, err := f.Create(ctx, v)
			if err != nil {
				return err
			}
		}
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			Model((*model.Subject)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		return err
	})
}
