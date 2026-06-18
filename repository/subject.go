package repository

import (
	"context"
	"database/sql"
	"errors"
	"sid/service/model"

	"github.com/uptrace/bun"
)

type Subject struct {
	db *bun.DB
}

func NewSubject(db *bun.DB) *Subject {
	Subject := Subject{db}
	return &Subject
}

func (r *Subject) SubjectCreate(ctx context.Context, data *model.Subject) (*model.Subject, error) {
	res := &model.Subject{}
	_, err := r.db.
		NewInsert().
		Model(data).
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *Subject) SubjectUpdate(ctx context.Context, data *model.Subject) (*model.Subject, error) {
	res := &model.Subject{}
	_, err := r.db.
		NewUpdate().
		Model(data).
		WherePK().
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *Subject) SubjectDelete(ctx context.Context, id *int64) (*model.Subject, error) {
	res := &model.Subject{}
	_, err := r.db.
		NewDelete().
		Model(&model.Subject{ID: *id}).
		WherePK().
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *Subject) SubjectAll(ctx context.Context) (*[]model.Subject, error) {
	res := &[]model.Subject{}
	err := r.db.
		NewSelect().
		Model((*model.Subject)(nil)).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (r *Subject) SubjectFindByID(ctx context.Context, id *int64) (*model.Subject, error) {
	res := &model.Subject{}
	err := r.db.
		NewSelect().
		Model((*model.Subject)(nil)).
		Where("id = ?", id).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return res, err
}
