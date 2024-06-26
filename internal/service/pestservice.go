/*
 * pestservice.go: Manages pest data for vineyards.
 * Provides CRUD and query functionality for pest observations.
 * Usage: Supports detailed management and analysis of pest data.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package service

import (
	"context"
	"errors"
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type PestService interface {
	CreatePestData(ctx context.Context, pest *model.PestData) error
	GetPestData(ctx context.Context, id int) (*model.PestData, error)
	UpdatePestData(ctx context.Context, pest *model.PestData) error
	DeletePestData(ctx context.Context, id int) error
	ListPestDataByVineyard(ctx context.Context, vineyardID int) ([]model.PestData, error)
	ListPestDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.PestData, error)
	FilterPestData(ctx context.Context, vineyardID int, pestType, severity string) ([]model.PestData, error)
}

type pestServiceImpl struct {
	db *db.DB
}

func NewPestService(db *db.DB) PestService {
	return &pestServiceImpl{db: db}
}

func (ps *pestServiceImpl) CreatePestData(ctx context.Context, pest *model.PestData) error {
	if pest == nil {
		return errors.New("cannot create nil pest data")
	}
	return ps.db.SavePestData(ctx, pest)
}

func (ps *pestServiceImpl) GetPestData(ctx context.Context, id int) (*model.PestData, error) {
	if id <= 0 {
		return nil, errors.New("invalid pest data ID")
	}
	return ps.db.GetPestData(ctx, id)
}

func (ps *pestServiceImpl) UpdatePestData(ctx context.Context, pest *model.PestData) error {
	if pest == nil {
		return errors.New("cannot update nil pest data")
	}
	if pest.ID == 0 {
		return errors.New("invalid pest data ID")
	}
	return ps.db.UpdatePestData(ctx, pest)
}

func (ps *pestServiceImpl) DeletePestData(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid pest data ID")
	}
	return ps.db.DeletePestData(ctx, id)
}

func (ps *pestServiceImpl) ListPestDataByVineyard(ctx context.Context, vineyardID int) ([]model.PestData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return ps.db.ListPestDataByVineyard(ctx, vineyardID)
}

func (ps *pestServiceImpl) ListPestDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.PestData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end date")
	}
	return ps.db.ListPestDataByDateRange(ctx, vineyardID, start, end)
}

func (ps *pestServiceImpl) FilterPestData(ctx context.Context, vineyardID int, pestType, severity string) ([]model.PestData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return ps.db.FilterPestData(ctx, vineyardID, pestType, severity)
}
