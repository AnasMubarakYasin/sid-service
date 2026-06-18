package repository

import (
	"context"
	"database/sql"
	"errors"
	"sid/service/model"

	"github.com/uptrace/bun"
)

type Experiment struct {
	db *bun.DB
}

func NewExperiment(db *bun.DB) *Experiment {
	Experiment := Experiment{db}
	return &Experiment
}

func (r *Experiment) ExperimentCreate(ctx context.Context, data *model.Experiment) (*model.Experiment, error) {
	res := &model.Experiment{}
	_, err := r.db.
		NewInsert().
		Model(data).
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *Experiment) ExperimentUpdate(ctx context.Context, data *model.Experiment) (*model.Experiment, error) {
	res := &model.Experiment{}
	_, err := r.db.
		NewUpdate().
		Model(data).
		WherePK().
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *Experiment) ExperimentDelete(ctx context.Context, id *int64) (*model.Experiment, error) {
	res := &model.Experiment{}
	_, err := r.db.
		NewDelete().
		Model(&model.Experiment{ID: *id}).
		WherePK().
		Returning("*").
		Exec(ctx, res)
	return res, err
}

func (r *Experiment) ExperimentDeleteAll(ctx context.Context) (*[]model.Experiment, error) {
	res := &[]model.Experiment{}
	_, err := r.db.
		NewTruncateTable().
		Model((*model.Experiment)(nil)).
		Exec(ctx, res)
	return res, err
}

func (r *Experiment) ExperimentAll(ctx context.Context) (*[]model.Experiment, error) {
	res := &[]model.Experiment{}
	err := r.db.
		NewSelect().
		Model((*model.Experiment)(nil)).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (r *Experiment) ExperimentFindByIDS(ctx context.Context, ids []int64) (*[]model.Experiment, error) {
	res := &[]model.Experiment{}
	err := r.db.
		NewSelect().
		Model((*model.Experiment)(nil)).
		Where("id IN (?)", bun.List(ids)).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		// Handle not found
		return nil, nil
	}
	return res, err
}

func (r *Experiment) ExperimentFindByID(ctx context.Context, id *int64) (*model.Experiment, error) {
	res := &model.Experiment{}
	err := r.db.
		NewSelect().
		Model((*model.Experiment)(nil)).
		Where("id = ?", id).
		Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (r *Experiment) ExperimentCount(ctx context.Context) (int, error) {
	count, err := r.db.
		NewSelect().
		Model((*model.Experiment)(nil)).
		Count(ctx)

	return count, err
}

func (r *Experiment) ExperimentExist(ctx context.Context, id *int64) (bool, error) {
	exists, err := r.db.
		NewSelect().
		Model((*Experiment)(nil)).
		Where("id = ?", id).
		Exists(ctx)
	if err != nil {
		panic(err)
	}
	return exists, err
}
