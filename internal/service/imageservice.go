/*
 * imageservice.go: Manages image data operations.
 * Handles CRUD operations and interfaces with storage solutions.
 * Usage: Provides methods to save, fetch, and manage images related to vineyards.
 * Author(s): Shannon Thompson
 * Created on: 04/11/2024
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

// ImageService defines the interface for image management, supporting CRUD operations and more.
type ImageService interface {
	SaveImage(ctx context.Context, image *model.Image, imageData io.Reader) error
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImage(ctx context.Context, image *model.Image) error
	DeleteImage(ctx context.Context, id int) error
	ListImagesByVineyard(ctx context.Context, vineyardID int) ([]model.Image, error)
	FindImagesByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.Image, error)
	GetRecentImages(ctx context.Context, vineyardID int, limit int) ([]model.Image, error)
}

// imageServiceImpl is the concrete implementation of ImageService using a database and storage service.
type imageServiceImpl struct {
	db      *db.DB
	storage *storage.StorageService
}

// NewImageService constructs a new ImageService given a database and a storage service instance.
func NewImageService(db *db.DB, storage *storage.StorageService) ImageService {
	return &imageServiceImpl{
		db:      db,
		storage: storage,
	}
}

// SaveImage handles the saving of a new image, both in the database and in cloud storage.
func (is *imageServiceImpl) SaveImage(ctx context.Context, image *model.Image, imageData io.Reader) error {
	if image == nil {
		return errors.New("cannot save nil image")
	}

	// Upload image data to cloud storage and retrieve the URL
	imageURL, err := is.storage.UploadFile(ctx, "vineyard_images/"+time.Now().Format("20060102_150405")+"_"+image.URL, imageData)
	if err != nil {
		return err
	}
	image.URL = imageURL // Update image URL with the URL from storage

	// Save image metadata in the database
	return is.db.SaveImage(ctx, image)
}

// GetImage retrieves an image by its ID from the database.
func (is *imageServiceImpl) GetImage(ctx context.Context, id int) (*model.Image, error) {
	if id <= 0 {
		return nil, errors.New("invalid image ID")
	}
	return is.db.GetImage(ctx, id)
}

// UpdateImage updates an existing image's metadata in the database.
func (is *imageServiceImpl) UpdateImage(ctx context.Context, image *model.Image) error {
	if image == nil {
		return errors.New("cannot update nil image")
	}
	if image.ID == 0 {
		return errors.New("invalid image ID")
	}
	return is.db.UpdateImage(ctx, image)
}

// DeleteImage removes an image's metadata from the database.
func (is *imageServiceImpl) DeleteImage(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid image ID")
	}
	return is.db.DeleteImage(ctx, id)
}

// ListImagesByVineyard retrieves all images associated with a specific vineyard.
func (is *imageServiceImpl) ListImagesByVineyard(ctx context.Context, vineyardID int) ([]model.Image, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return is.db.ListImagesByVineyard(ctx, vineyardID)
}

// FindImagesByDateRange searches for images within a specific date range and vineyard.
func (is *imageServiceImpl) FindImagesByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.Image, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end date")
	}
	return is.db.FindImagesByDateRange(ctx, vineyardID, start, end)
}

// GetRecentImages fetches the most recent images up to a specified limit for a vineyard.
func (is *imageServiceImpl) GetRecentImages(ctx context.Context, vineyardID int, limit int) ([]model.Image, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if limit <= 0 {
		return nil, errors.New("limit must be a positive number")
	}
	return is.db.GetRecentImages(ctx, vineyardID, limit)
}
