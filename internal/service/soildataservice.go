/*
 * File: soildataservice.go
 * Description: Manages soil data interactions, providing services to retrieve and update soil health
 *              information from various data sources configured within the application.
 * Usage:
 *   - Facilitates operations such as GetSoilData and UpdateSoilData to manage soil information linked to vineyards.
 *   - Interacts with external soil data APIs or local databases to fetch and store soil data.
 * Dependencies:
 *   - db.go for executing SQL queries related to soil data.
 *   - client.go or similar modules for fetching soil data from external APIs.
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

type SoilDataService interface {
	CreateSoilData(ctx context.Context, soilData *model.SoilData) error
	GetSoilData(ctx context.Context, id int) (*model.SoilData, error)
	UpdateSoilData(ctx context.Context, soilData *model.SoilData) error
	DeleteSoilData(ctx context.Context, id int) error
	ListSoilData(ctx context.Context, vineyardID int) ([]model.SoilData, error)
	ListSoilDataByVineyard(ctx context.Context, vineyardID int) ([]model.SoilData, error)
	ListSoilDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SoilData, error)
}

type soilDataServiceImpl struct {
	db *db.DB
}

func NewSoilDataService(db *db.DB) SoilDataService {
	return &soilDataServiceImpl{db: db}
}

func (sds *soilDataServiceImpl) CreateSoilData(ctx context.Context, soilData *model.SoilData) error {
	if soilData == nil {
		return errors.New("cannot create nil soil data")
	}
	return sds.db.SaveSoilData(ctx, soilData)
}

func (sds *soilDataServiceImpl) GetSoilData(ctx context.Context, id int) (*model.SoilData, error) {
	if id <= 0 {
		return nil, errors.New("invalid soil data ID")
	}
	return sds.db.GetSoilData(ctx, id)
}

func (sds *soilDataServiceImpl) UpdateSoilData(ctx context.Context, soilData *model.SoilData) error {
	if soilData == nil {
		return errors.New("cannot update nil soil data")
	}
	if soilData.ID == 0 {
		return errors.New("invalid soil data ID")
	}
	return sds.db.UpdateSoilData(ctx, soilData)
}

func (sds *soilDataServiceImpl) DeleteSoilData(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid soil data ID")
	}
	return sds.db.DeleteSoilData(ctx, id)
}

func (sds *soilDataServiceImpl) ListSoilData(ctx context.Context, vineyardID int) ([]model.SoilData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return sds.db.ListSoilDataForVineyard(ctx, vineyardID)
}

func (sds *soilDataServiceImpl) ListSoilDataByVineyard(ctx context.Context, vineyardID int) ([]model.SoilData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return sds.db.ListSoilDataForVineyard(ctx, vineyardID)
}

func (sds *soilDataServiceImpl) ListSoilDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SoilData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return sds.db.ListSoilDataByDateRange(ctx, vineyardID, start, end)
}
