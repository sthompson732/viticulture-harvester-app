/*
 * satelliteservice.go: Manages satellite imagery data for vineyards.
 * Provides concurrent operations for storing, updating, and retrieving high-resolution imagery efficiently.
 * Usage: Manages satellite data integrations, storage operations, and concurrent data queries.
 * satelliteservice.go: Manages satellite imagery data for vineyards.
 * Provides concurrent operations for storing, updating, and retrieving high-resolution imagery efficiently.
 * Usage: Manages satellite data integrations, storage operations, and concurrent data queries.
 * Author(s): Shannon Thompson
 * Created on: 04/12/2024
 * Created on: 04/12/2024
 */
package service

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync"
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
	"github.com/sthompson732/viticulture-harvester-app/internal/storage"
)

type SatelliteService interface {
	SaveSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error
	GetSatelliteData(ctx context.Context, id int) (*model.SatelliteData, error)
	UpdateSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error
	DeleteSatelliteData(ctx context.Context, id int) error
	ListSatelliteDataByVineyard(ctx context.Context, vineyardID int) ([]model.SatelliteData, error)
	ListSatelliteImageryByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SatelliteData, error)
	ConcurrentSaveSatelliteData(ctx context.Context, datas []*model.SatelliteData, imageDatas []io.Reader) error
	ConcurrentSaveSatelliteData(ctx context.Context, datas []*model.SatelliteData, imageDatas []io.Reader) error
}

type satelliteServiceImpl struct {
	db      *db.DB
	storage *storage.StorageService
}

func NewSatelliteService(db *db.DB, storage *storage.StorageService) SatelliteService {
	return &satelliteServiceImpl{db: db, storage: storage}
}

func (s *satelliteServiceImpl) SaveSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error {
	if data == nil {
		return errors.New("cannot save nil satellite data")
	}

	// Upload image data to cloud storage and retrieve the URL
	imageURL, err := s.storage.UploadFile(ctx, "satellite_images/"+data.ImageURL, imageData)
	if err != nil {
		return err
	}
	data.ImageURL = imageURL // Update image URL with the URL from storage

	// Save satellite data metadata in the database
	return s.db.SaveSatelliteImageryMetadata(ctx, data, data.VineyardID)
}

func (s *satelliteServiceImpl) GetSatelliteData(ctx context.Context, id int) (*model.SatelliteData, error) {
	if id <= 0 {
		return nil, errors.New("invalid satellite data ID")
	}
	return s.db.GetSatelliteImagery(ctx, id)
}

func (s *satelliteServiceImpl) UpdateSatelliteData(ctx context.Context, data *model.SatelliteData, imageData io.Reader) error {
	if data == nil || data.ID == 0 {
		return errors.New("invalid satellite data")
	if data == nil || data.ID == 0 {
		return errors.New("invalid satellite data")
	}
	imageURL, err := s.storage.UploadFile(ctx, "satellite_images/"+data.ImageURL, imageData)
	if err != nil {
		return err
	}
	data.ImageURL = imageURL
	return s.db.UpdateSatelliteImagery(ctx, data)
}

func (s *satelliteServiceImpl) DeleteSatelliteData(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid satellite data ID")
	}
	return s.db.DeleteSatelliteImagery(ctx, id)
}

func (s *satelliteServiceImpl) ListSatelliteDataByVineyard(ctx context.Context, vineyardID int) ([]model.SatelliteData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return s.db.ListSatelliteImageryByVineyard(ctx, vineyardID)
}

func (s *satelliteServiceImpl) ListSatelliteImageryByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SatelliteData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end date")
	}
	return s.db.ListSatelliteImageryByDateRange(ctx, vineyardID, start, end)
}

func (s *satelliteServiceImpl) ConcurrentSaveSatelliteData(ctx context.Context, datas []*model.SatelliteData, imageDatas []io.Reader) error {
	if len(datas) != len(imageDatas) {
		return errors.New("data and image slices must be of the same length")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(datas))

	for i, data := range datas {
		wg.Add(1)
		go func(data *model.SatelliteData, imageData io.Reader) {
			defer wg.Done()
			if err := s.SaveSatelliteData(ctx, data, imageData); err != nil {
				errChan <- err
			}
		}(data, imageDatas[i])
	}

	wg.Wait()
	close(errChan)

	// Check if there were any errors
	for err := range errChan {
		if err != nil {
			return err // Return the first encountered error
		}
	}

	return nil
}
