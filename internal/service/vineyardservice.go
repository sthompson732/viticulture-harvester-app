package service

import (
	"context"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
	"github.com/sthompson732/viticulture-harvester-app/internal/storage"
)

type VineyardService interface {
	CreateVineyard(ctx context.Context, vineyard *model.Vineyard) error
	GetVineyard(ctx context.Context, id string) (*model.Vineyard, error)
	UpdateVineyard(ctx context.Context, id string, vineyard *model.Vineyard) error
	DeleteVineyard(ctx context.Context, id string) error
	ListVineyards(ctx context.Context) ([]model.Vineyard, error)
}

// Implementation of VineyardService
type vineyardServiceImpl struct {
	db      *db.DB
	storage *storage.StorageService
}

func NewVineyardService(db *db.DB, storage *storage.StorageService) VineyardService {
	return &vineyardServiceImpl{db: db, storage: storage}
}

func (vs *vineyardServiceImpl) CreateVineyard(ctx context.Context, vineyard *model.Vineyard) error {
	return vs.db.SaveVineyard(ctx, vineyard)
}

func (vs *vineyardServiceImpl) GetVineyard(ctx context.Context, id string) (*model.Vineyard, error) {
	return vs.db.GetVineyard(ctx, id)
}

func (vs *vineyardServiceImpl) UpdateVineyard(ctx context.Context, id string, vineyard *model.Vineyard) error {
	return vs.db.UpdateVineyard(ctx, id, vineyard)
}

func (vs *vineyardServiceImpl) DeleteVineyard(ctx context.Context, id string) error {
	return vs.db.DeleteVineyard(ctx, id)
}

func (vs *vineyardServiceImpl) ListVineyards(ctx context.Context) ([]model.Vineyard, error) {
	return vs.db.ListVineyards(ctx)
}
