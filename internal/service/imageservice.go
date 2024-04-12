/*
 * File: imageservice.go
 * Description: Provides services for managing image data related to vineyards. Handles CRUD
 *              operations for images, interfacing with both the database and external storage solutions.
 * Usage:
 *   - Offers methods like SaveImage, GetImage, DeleteImage, and ListImages to manage image records.
 *   - Utilizes storage interfaces to save and retrieve image files, typically from cloud storage.
 * Dependencies:
 *   - db.go for database operations related to image data.
 *   - storage.go for interfacing with external file storage systems (e.g., AWS S3, Google Cloud Storage).
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

type ImageService interface {
	SaveImage(ctx context.Context, image *model.Image) error
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImage(ctx context.Context, image *model.Image) error
	DeleteImage(ctx context.Context, id int) error
	ListImagesByVineyard(ctx context.Context, vineyardID int) ([]model.Image, error)
	FindImagesByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.Image, error)
	GetRecentImages(ctx context.Context, vineyardID int, limit int) ([]model.Image, error)
}

type imageServiceImpl struct {
	db *db.DB
}

func NewImageService(db *db.DB) ImageService {
	return &imageServiceImpl{db: db}
}

func (is *imageServiceImpl) SaveImage(ctx context.Context, image *model.Image) error {
	if image == nil {
		return errors.New("cannot save nil image")
	}
	return is.db.SaveImage(ctx, image)
}

func (is *imageServiceImpl) GetImage(ctx context.Context, id int) (*model.Image, error) {
	if id <= 0 {
		return nil, errors.New("invalid image ID")
	}
	return is.db.GetImage(ctx, id)
}

func (is *imageServiceImpl) UpdateImage(ctx context.Context, image *model.Image) error {
	if image == nil {
		return errors.New("cannot update nil image")
	}
	if image.ID == 0 {
		return errors.New("invalid image ID")
	}
	return is.db.UpdateImage(ctx, image)
}

func (is *imageServiceImpl) DeleteImage(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid image ID")
	}
	return is.db.DeleteImage(ctx, id)
}

func (is *imageServiceImpl) ListImagesByVineyard(ctx context.Context, vineyardID int) ([]model.Image, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return is.db.ListImagesForVineyard(ctx, vineyardID)
}

func (is *imageServiceImpl) FindImagesByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.Image, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end date")
	}
	return is.db.FindImagesByDateRange(ctx, vineyardID, start, end)
}

func (is *imageServiceImpl) GetRecentImages(ctx context.Context, vineyardID int, limit int) ([]model.Image, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if limit <= 0 {
		return nil, errors.New("limit must be a positive number")
	}
	return is.db.GetRecentImages(ctx, vineyardID, limit)
}
