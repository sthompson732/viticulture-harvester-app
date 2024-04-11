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
	"strconv" // For converting string IDs to integers

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type SoilDataService interface {
	UpdateSoilData(ctx context.Context, vineyardID string, soilData *model.SoilData) error
	GetSoilData(ctx context.Context, vineyardID string) (*model.SoilData, error)
}

type soilDataServiceImpl struct {
	db *db.DB
}

func NewSoilDataService(db *db.DB) SoilDataService {
	return &soilDataServiceImpl{db: db}
}

func (sds *soilDataServiceImpl) UpdateSoilData(ctx context.Context, vineyardID string, soilData *model.SoilData) error {
	intID, err := strconv.Atoi(vineyardID) // Convert id from string to int
	if err != nil {
		return err // Proper error handling for ID conversion
	}
	return sds.db.UpdateSoilData(ctx, intID, soilData)
}

func (sds *soilDataServiceImpl) GetSoilData(ctx context.Context, vineyardID string) (*model.SoilData, error) {
	intID, err := strconv.Atoi(vineyardID) // Convert id from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return sds.db.GetSoilDataForVineyard(ctx, intID)
}
