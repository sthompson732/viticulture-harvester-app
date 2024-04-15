/*
 * storage.go: Manages file operations in cloud storage.
 * Implements file saving, retrieval, and deletion within Google Cloud Storage.
 * Usage: Supports imageservice and satelliteservice with file management capabilities.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// StorageService encapsulates the Google Cloud Storage client and related operations.
type StorageService struct {
	Client     *storage.Client
	BucketName string
}

// NewStorageService initializes a new storage service with the provided Google Cloud Storage bucket.
func NewStorageService(ctx context.Context, bucketName string) (*StorageService, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}
	return &StorageService{
		Client:     client,
		BucketName: bucketName,
	}, nil
}

// UploadImage uploads an image to the cloud storage and returns its URL.
func (s *StorageService) UploadImage(ctx context.Context, fileName string, imageData io.Reader) (string, error) {
	bucket := s.Client.Bucket(s.BucketName)
	obj := bucket.Object(fileName)

	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, imageData); err != nil {
		return "", fmt.Errorf("failed to write image to bucket: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize image upload: %w", err)
	}

	// Set the image to be publicly readable.
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set image public: %w", err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get uploaded image attributes: %w", err)
	}

	return attrs.MediaLink, nil
}

// ListImages retrieves a list of image URLs from the cloud storage.
func (s *StorageService) ListImages(ctx context.Context) ([]string, error) {
	var urls []string
	it := s.Client.Bucket(s.BucketName).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list images: %w", err)
		}
		urls = append(urls, attrs.MediaLink)
	}
	return urls, nil
}
