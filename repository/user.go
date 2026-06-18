package repository

import (
	"context"
	"database/sql"
	"errors"
	"sid/service/model"

	"github.com/uptrace/bun"
)

// type UserCreate struct {
// 	AccountID int64

// 	Photo    string
// 	Name     string
// 	Email    string
// 	Settings map[string]any
// 	Status   model.UserStatus
// }

type User struct {
	db *bun.DB
}

func NewUser(db *bun.DB) *User {
	user := User{db}
	return &user
}

// func (r *User) ctx context.Context, UserCreate(data *UserCreate) (*model.User, error) {
// 	res := &model.User{}
// 	_, err := r.db.
// 		NewInsert().
// 		Model(&model.User{
// 			AccountID: data.AccountID,
// 			Photo:     data.Photo,
// 			Name:      data.Name,
// 			Email:     data.Email,
// 			Settings:  data.Settings,
// 			Status:    data.Status,
// 		}
// 	return res, err
// }

func (r *User) UserCreate(ctx context.Context, data *model.User) (*model.User, error) {
	res := &model.User{}
	_, err := r.db.
		NewInsert().
		Model(data).
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *User) UserUpdate(ctx context.Context, id int64, data *model.User) (*model.User, error) {
	res := &model.User{}
	_, err := r.db.
		NewUpdate().
		Model(data).
		WherePK().
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *User) UserDelete(ctx context.Context, id int64) (*model.User, error) {
	res := &model.User{}
	_, err := r.db.
		NewDelete().
		Model(&model.User{ID: id}).
		WherePK().
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *User) UserFind(ctx context.Context, data *model.User) (*[]model.User, error) {
	res := &[]model.User{}
	err := r.db.
		NewSelect().
		Model(data).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		// Handle not found
		return nil, nil
	}
	return res, err
}

func (r *User) UserFindByID(ctx context.Context, id int64) (*model.User, error) {
	res := &model.User{}
	err := r.db.
		NewSelect().
		Model((*model.User)(nil)).
		Where("id = ?", id).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		// Handle not found
		return nil, nil
	}
	return res, err
}

func (r *User) UserFindByIDS(ctx context.Context, ids []int64) (*[]model.User, error) {
	res := &[]model.User{}
	err := r.db.
		NewSelect().
		Model((*model.User)(nil)).
		Where("id IN (?)", bun.List(ids)).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		// Handle not found
		return nil, nil
	}
	return res, err
}

func (r *User) UserFindByName(ctx context.Context, name string) (*model.User, error) {
	res := &model.User{}
	err := r.db.
		NewSelect().
		Model((*model.User)(nil)).
		Where("name = ?", name).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		// Handle not found
		return nil, nil
	}
	return res, err
}
