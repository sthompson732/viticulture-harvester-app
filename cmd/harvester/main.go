/*
 * main.go: Entry point of the Viticulture Data Harvester.
 * Initializes services and dynamically schedules data fetching tasks.
 * Usage: Sets up services, routing, and triggers data ingestion based on configuration.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package main

import (
	"context"
	"log"
	"os"

	"github.com/sthompson732/viticulture-harvester-app/internal/api"
	"github.com/sthompson732/viticulture-harvester-app/internal/config"
	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/scheduler"
	"github.com/sthompson732/viticulture-harvester-app/internal/server"
	"github.com/sthompson732/viticulture-harvester-app/internal/service"
	"github.com/sthompson732/viticulture-harvester-app/internal/storage"
)

func main() {
	ctx := context.Background()

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
	storageService, err := storage.NewStorageService(ctx, cfg.CloudStorage.BucketName, cfg.CloudStorage.CredentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	// Initialize data services
	vineyardService := service.NewVineyardService(database)
	imageService := service.NewImageService(database, storageService)
	soilDataService := service.NewSoilDataService(database)
	pestService := service.NewPestService(database)
	weatherService := service.NewWeatherService(database)
	satelliteService := service.NewSatelliteService(database, storageService)

	// Set up the API router
	router := api.NewRouter(vineyardService, imageService, soilDataService, pestService, weatherService, satelliteService)

	// Initialize and start the server
	srv := server.NewServer(router)
	if err := srv.Start(cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Initialize Scheduler Client and set up jobs
	schedClient, err := scheduler.NewSchedulerClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create scheduler client: %v", err)
	}

	// Dynamically schedule jobs based on data source configurations
	if err := schedClient.SetupJobs(ctx); err != nil {
		log.Fatalf("Failed to set up scheduler jobs: %v", err)
	}
}
