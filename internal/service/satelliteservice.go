/*
 * satelliteservice.go: Handles satellite imagery data for vineyards.
 * Provides functions for storing, updating, and retrieving high-resolution imagery.
 * Usage: Used to manage satellite data integrations and queries.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package service

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
	"github.com/sthompson732/viticulture-harvester-app/internal/storage"
)

// SatelliteService defines the interface for satellite data management, including CRUD operations and listings by criteria.
type SatelliteService interface {
	SaveSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error
	GetSatelliteData(ctx context.Context, id int) (*model.SatelliteData, error)
	UpdateSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error
	DeleteSatelliteData(ctx context.Context, id int) error
	ListSatelliteDataByVineyard(ctx context.Context, vineyardID int) ([]model.SatelliteData, error)
	ListSatelliteImageryByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SatelliteData, error)
}

type satelliteServiceImpl struct {
	db      *db.DB
	storage *storage.StorageService
}

// NewSatelliteService creates a new instance of a satellite service that uses specific database and storage service implementations.
func NewSatelliteService(db *db.DB, storage *storage.StorageService) SatelliteService {
	return &satelliteServiceImpl{db: db, storage: storage}
}

// SaveSatelliteData handles the saving of new satellite data along with uploading the associated image to cloud storage.
func (s *satelliteServiceImpl) SaveSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error {
	if data == nil {
		return errors.New("cannot save nil satellite data")
	}
	if imageData != nil {
		imageURL, err := s.storage.UploadImage(ctx, "satellite_images/"+data.ImageURL, imageData)
		if err != nil {
			return err
		}
		data.ImageURL = imageURL
	}
	return s.db.SaveSatelliteImageryMetadata(ctx, data, data.VineyardID)
}

// GetSatelliteData retrieves satellite data by its ID.
func (s *satelliteServiceImpl) GetSatelliteData(ctx context.Context, id int) (*model.SatelliteData, error) {
	if id <= 0 {
		return nil, errors.New("invalid satellite data ID")
	}
	return s.db.GetSatelliteImagery(ctx, id)
}

// UpdateSatelliteData updates the metadata for an existing set of satellite data; can also update the associated image.
func (s *satelliteServiceImpl) UpdateSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error {
	if data == nil {
		return errors.New("cannot update nil satellite data")
	}
	if data.ID == 0 {
		return errors.New("invalid satellite data ID")
	}
	if imageData != nil {
		imageURL, err := s.storage.UploadImage(ctx, "satellite_images/"+data.ImageURL, imageData)
		if err != nil {
			return err
		}
		data.ImageURL = imageURL
	}
	return s.db.UpdateSatelliteImagery(ctx, data)
}

// DeleteSatelliteData removes satellite data from the database by its ID.
func (s *satelliteServiceImpl) DeleteSatelliteData(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid satellite data ID")
	}
	return s.db.DeleteSatelliteImagery(ctx, id)
}

// ListSatelliteDataByVineyard lists all satellite data associated with a specific vineyard.
func (s *satelliteServiceImpl) ListSatelliteDataByVineyard(ctx context.Context, vineyardID int) ([]model.SatelliteData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return s.db.ListSatelliteImageryByVineyard(ctx, vineyardID)
}

// ListSatelliteImageryByDateRange lists satellite data for a specific vineyard within a given date range.
func (s *satelliteServiceImpl) ListSatelliteImageryByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SatelliteData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end date")
	}
	return s.db.ListSatelliteImageryByDateRange(ctx, vineyardID, start, end)
}
