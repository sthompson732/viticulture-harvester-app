/*
 * File: vineyardservice.go
 * Description: Handles business logic related to vineyard management. Provides methods for creating,
 *              retrieving, updating, and deleting vineyard entries, as well as interfacing with related data like
 *              images and soil data.
 * Usage:
 *   - Implements business processes associated with vineyards, such as registering new vineyards or updating
 *     existing entries.
 *   - Coordinates with db.go to perform database transactions and with other services for related data.
 * Dependencies:
 *   - Relies on db.go for database interactions.
 *   - Integrates with imageservice.go and soildataservice.go for comprehensive data management.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package service

import (
	"context"
	"strconv" // For converting string IDs to integers if IDs are integers in the database

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
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return vs.db.GetVineyard(ctx, intID)
}

func (vs *vineyardServiceImpl) UpdateVineyard(ctx context.Context, id string, vineyard *model.Vineyard) error {
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return err // Proper error handling for ID conversion
	}
	return vs.db.UpdateVineyard(ctx, intID, vineyard)
}

func (vs *vineyardServiceImpl) DeleteVineyard(ctx context.Context, id string) error {
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return err // Proper error handling for ID conversion
	}
	return vs.db.DeleteVineyard(ctx, intID)
}

func (vs *vineyardServiceImpl) ListVineyards(ctx context.Context) ([]model.Vineyard, error) {
	return vs.db.ListVineyards(ctx)
}
