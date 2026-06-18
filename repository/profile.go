package repository

import (
	"context"
	"sid/service/model"

	"github.com/uptrace/bun"
)

type Profile struct {
	db *bun.DB
}

func NewProfile(db *bun.DB) *Profile {
	Profile := Profile{}
	return &Profile
}

func (r *Profile) AccountCreate(ctx context.Context, data *model.User) (*model.User, error) {
	res := &model.User{}
	_, err := r.db.
		NewInsert().
		Model(data).
		Exec(ctx, res)
	return res, err
}

func (r *Profile) AccountUpdate(ctx context.Context, id int64, data *model.User) (*model.User, error) {
	res := &model.User{}
	_, err := r.db.
		NewUpdate().
		Model(data).
		WherePK().
		Exec(ctx, res)
	return res, err
}

func (r *Profile) AccountDelete(ctx context.Context, id int64) (*model.User, error) {
	res := &model.User{}
	_, err := r.db.
		NewDelete().
		Model(&model.User{ID: id}).
		WherePK().
		Exec(ctx, res)
	return res, err
}

func (r *Profile) AccountFind(ctx context.Context, data *model.User) (*[]model.User, error) {
	res := &[]model.User{}
	err := r.db.
		NewSelect().
		Model(data).
		Scan(ctx, res)
	return res, err
}

func (r *Profile) AccountFindByID(ctx context.Context, id int64) (*model.User, error) {
	res := &model.User{}
	err := r.db.
		NewSelect().
		Model((*model.User)(nil)).
		Where("id = ?", id).
		Scan(ctx, res)
	return res, err
}
