/*
 * File: main.go
 * Description: Entry point of the Viniculture Data Harvester application. This file initializes
 *              all necessary services, sets up the HTTP router, and starts the server.
 * Usage:
 *   - Initializes database connections, configures services, and prepares the HTTP server.
 *   - Routes are set up using the router.go configurations, and the server is started on a specified port.
 * Dependencies:
 *   - server.go for starting the server.
 *   - router.go for HTTP routing configurations.
 *   - service files like vineyardservice.go for business logic.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package main

import (
	"context"
	"log"
	"os"

	"github.com/sthompson732/viticulture-harvester-app/internal/api"
	"github.com/sthompson732/viticulture-harvester-app/internal/clients"
	"github.com/sthompson732/viticulture-harvester-app/internal/config"
	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/scheduler"
	"github.com/sthompson732/viticulture-harvester-app/internal/server"
	"github.com/sthompson732/viticulture-harvester-app/internal/service"
	"github.com/sthompson732/viticulture-harvester-app/internal/storage"
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
	storageService, err := storage.NewStorageService(ctx, cfg.CloudStorage.BucketName, cfg.CloudStorage.CredentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	// Initialize services
	vineyardService := service.NewVineyardService(database)
	imageService := service.NewImageService(database, storageService)
	soilDataService := service.NewSoilDataService(database)
	pestService := service.NewPestService(database)
	weatherService := service.NewWeatherService(database)
	satelliteService := service.NewSatelliteService(database, storageService)

	// Set up the router
	router := api.NewRouter(vineyardService, imageService, soilDataService, pestService, weatherService, satelliteService)

	// Initialize and start the server
	srv := server.NewServer(router)
	if err := srv.Start(cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Initialize satellite and soil clients
	satelliteClient := clients.NewSatelliteClient(cfg)
	soilClient := clients.NewSoilClient(cfg)

	// Initialize Scheduler Client and set up jobs
	schedClient, err := scheduler.NewSchedulerClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create scheduler client: %v", err)
	}

	if err := schedClient.InitializeJobs(ctx); err != nil {
		log.Fatalf("Failed to initialize scheduler jobs: %v", err)
	}
}
