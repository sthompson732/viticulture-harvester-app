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
	"strconv" // For converting string IDs to integers

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
	"github.com/sthompson732/viticulture-harvester-app/internal/storage"
)

type ImageService interface {
	ListImages(ctx context.Context, vineyardID string) ([]model.Image, error)
	SaveImage(ctx context.Context, image *model.Image) error
	GetImage(ctx context.Context, id string) (*model.Image, error)
	DeleteImage(ctx context.Context, id string) error
}

type imageServiceImpl struct {
	db      *db.DB
	storage *storage.StorageService
}

func NewImageService(db *db.DB, storage *storage.StorageService) ImageService {
	return &imageServiceImpl{db: db, storage: storage}
}

func (is *imageServiceImpl) ListImages(ctx context.Context, vineyardID string) ([]model.Image, error) {
	intID, err := strconv.Atoi(vineyardID) // Convert vineyardID from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return is.db.GetSatelliteImageryForVineyard(ctx, intID)
}

func (is *imageServiceImpl) SaveImage(ctx context.Context, image *model.Image) error {
	return is.db.SaveImage(ctx, image)
}

func (is *imageServiceImpl) GetImage(ctx context.Context, id string) (*model.Image, error) {
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return nil, err // Proper error handling for ID conversion
	}
	return is.db.GetImage(ctx, intID)
}

func (is *imageServiceImpl) DeleteImage(ctx context.Context, id string) error {
	intID, err := strconv.Atoi(id) // Convert id from string to int
	if err != nil {
		return err // Proper error handling for ID conversion
	}
	return is.db.DeleteImage(ctx, intID)
}
