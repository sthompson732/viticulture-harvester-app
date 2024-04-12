/*
 * File: satelliteservice.go
 * Description: Manages satellite imagery data interactions, providing services to retrieve, update, and delete satellite
 *              imagery information related to vineyards.
 * Usage:
 *   - Facilitates operations such as storing new satellite images, updating metadata, deleting entries, and retrieving
 *     images based on various criteria such as vineyard ID or specific date ranges.
 * Dependencies:
 *   - Relies on db.go for executing SQL queries related to satellite imagery data.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package service

import (
	"context"
	"strconv" // For converting string IDs to integers
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type SatelliteService interface {
	SaveSatelliteData(ctx context.Context, data *model.SatelliteData) error
	GetSatelliteData(ctx context.Context, id string) (*model.SatelliteData, error)
	UpdateSatelliteData(ctx context.Context, data *model.SatelliteData) error
	DeleteSatelliteData(ctx context.Context, id string) error
	ListSatelliteDataByVineyard(ctx context.Context, vineyardID string) ([]model.SatelliteData, error)
	ListSatelliteDataByDateRange(ctx context.Context, vineyardID string, start, end time.Time) ([]model.SatelliteData, error)
}

type satelliteServiceImpl struct {
	db *db.DB
}

func NewSatelliteService(db *db.DB) SatelliteService {
	return &satelliteServiceImpl{db: db}
}

func (s *satelliteServiceImpl) SaveSatelliteData(ctx context.Context, data *model.SatelliteData) error {
	return s.db.SaveSatelliteImageryMetadata(ctx, data, data.VineyardID)
}

func (s *satelliteServiceImpl) GetSatelliteData(ctx context.Context, id string) (*model.SatelliteData, error) {
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return s.db.GetSatelliteImagery(ctx, intID)
}

func (s *satelliteServiceImpl) UpdateSatelliteData(ctx context.Context, data *model.SatelliteData) error {
	return s.db.UpdateSatelliteImagery(ctx, data)
}

func (s *satelliteServiceImpl) DeleteSatelliteData(ctx context.Context, id string) error {
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return err // Proper error handling for ID conversion
	}
	return s.db.DeleteSatelliteImagery(ctx, intID)
}

func (s *satelliteServiceImpl) ListSatelliteDataByVineyard(ctx context.Context, vineyardID string) ([]model.SatelliteData, error) {
	intID, err := strconv.Atoi(vineyardID) // Convert vineyardID from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return s.db.ListSatelliteImageryByVineyard(ctx, intID)
}

func (s *satelliteServiceImpl) ListSatelliteDataByDateRange(ctx context.Context, vineyardID string, start, end time.Time) ([]model.SatelliteData, error) {
	intID, err := strconv.Atoi(vineyardID) // Convert vineyardID from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return s.db.ListSatelliteImageryByDateRange(ctx, intID, start, end)
}
