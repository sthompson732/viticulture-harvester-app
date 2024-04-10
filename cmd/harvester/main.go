package main

import (
	"context"
	"log"
	"os"
	"viticulture-harvester-app/internal/clients"
	"viticulture-harvester-app/internal/config"
	"viticulture-harvester-app/internal/db"
	"viticulture-harvester-app/internal/server"
	"viticulture-harvester-app/internal/storage"
)

func main() {
	// Load configuration from file
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database
	database, err := db.NewDB(cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Initialize the storage service
	ctx := context.Background() // Depending on our needs, we might want to use a more specific context
	storageClient, err := storage.NewStorageService(ctx, cfg.CloudStorage.BucketName)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	// Initialize satellite and soil clients
	satelliteClient := clients.NewSatelliteClient(cfg)
	soilClient := clients.NewSoilClient(cfg)

	// Initialize and start the server
	srv := server.NewServer(database, storageClient, satelliteClient, soilClient)
	if err := srv.Start(cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
