/*
 * storage.go: Manages file operations in cloud storage.
 * Implements file saving, retrieval, and deletion within Google Cloud Storage.
 * Usage: Supports imageservice and satelliteservice with file management capabilities.
 * Author(s): Shannon Thompson
 * Created on: 04/12/2024
 */

package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// StorageService encapsulates the Google Cloud Storage client and related operations.
type StorageService struct {
	Client     *storage.Client
	BucketName string
}

// NewStorageService initializes a new storage service with the provided Google Cloud Storage bucket.
func NewStorageService(ctx context.Context, bucketName, credentialsPath string) (*StorageService, error) {
	// If credentials path is provided, use it to authenticate the client.
	var client *storage.Client
	var err error
	if credentialsPath != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	} else {
		// Otherwise, use the default credentials.
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &StorageService{
		Client:     client,
		BucketName: bucketName,
	}, nil
}

// UploadFile uploads any file to the cloud storage and returns its URL.
func (s *StorageService) UploadFile(ctx context.Context, filePath string, fileData io.Reader) (string, error) {
	bucket := s.Client.Bucket(s.BucketName)
	obj := bucket.Object(filePath)

	w := obj.NewWriter(ctx)

	// Create a buffer to store a snippet of the file data for content type detection
	buf := make([]byte, 512) // 512 bytes should be enough for content type detection
	n, err := fileData.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file for content type detection: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(buf)
	w.ContentType = contentType

	// Write the buffer read for detection
	if _, err := w.Write(buf[:n]); err != nil {
		return "", fmt.Errorf("failed to write initial data to bucket: %w", err)
	}

	// Continue writing the rest of the file
	if _, err := io.Copy(w, fileData); err != nil {
		return "", fmt.Errorf("failed to write file to bucket: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize file upload: %w", err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set file public: %w", err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get uploaded file attributes: %w", err)
	}

	return attrs.MediaLink, nil
}

// DeleteFile deletes a file from cloud storage.
func (s *StorageService) DeleteFile(ctx context.Context, filePath string) error {
	bucket := s.Client.Bucket(s.BucketName)
	obj := bucket.Object(filePath)

	err := obj.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFileMetadata retrieves metadata of a file stored in cloud storage.
func (s *StorageService) GetFileMetadata(ctx context.Context, filePath string) (*storage.ObjectAttrs, error) {
	bucket := s.Client.Bucket(s.BucketName)
	obj := bucket.Object(filePath)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file attributes: %w", err)
	}
	return attrs, nil
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
