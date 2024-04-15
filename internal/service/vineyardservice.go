/*
 * vineyardservice.go: Business logic for vineyard management.
 * Provides functions for CRUD operations on vineyards.
 * Usage: Interacts with database for vineyard data manipulation.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package service

import (
	"context"
	"errors"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type VineyardService interface {
	CreateVineyard(ctx context.Context, vineyard *model.Vineyard) error
	GetVineyard(ctx context.Context, id int) (*model.Vineyard, error)
	UpdateVineyard(ctx context.Context, vineyard *model.Vineyard) error
	DeleteVineyard(ctx context.Context, id int) error
	ListVineyards(ctx context.Context) ([]model.Vineyard, error)
	GetVineyardWithEnvironmentalData(ctx context.Context, id int) (*model.Vineyard, error)
}

type vineyardServiceImpl struct {
	db *db.DB
}

func NewVineyardService(db *db.DB) VineyardService {
	return &vineyardServiceImpl{db: db}
}

func (vs *vineyardServiceImpl) CreateVineyard(ctx context.Context, vineyard *model.Vineyard) error {
	if vineyard == nil {
		return errors.New("cannot create a nil vineyard")
	}
	return vs.db.SaveVineyard(ctx, vineyard)
}

func (vs *vineyardServiceImpl) GetVineyard(ctx context.Context, id int) (*model.Vineyard, error) {
	if id <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return vs.db.GetVineyard(ctx, id)
}

func (vs *vineyardServiceImpl) UpdateVineyard(ctx context.Context, vineyard *model.Vineyard) error {
	if vineyard == nil {
		return errors.New("cannot update a nil vineyard")
	}
	if vineyard.ID == 0 {
		return errors.New("invalid vineyard ID")
	}
	return vs.db.UpdateVineyard(ctx, vineyard)
}

func (vs *vineyardServiceImpl) DeleteVineyard(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid vineyard ID")
	}
	return vs.db.DeleteVineyard(ctx, id)
}

func (vs *vineyardServiceImpl) ListVineyards(ctx context.Context) ([]model.Vineyard, error) {
	return vs.db.ListVineyards(ctx)
}

func (vs *vineyardServiceImpl) GetVineyardWithEnvironmentalData(ctx context.Context, id int) (*model.Vineyard, error) {
	if id <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return vs.db.GetVineyardWithEnvironmentalData(ctx, id)
}
