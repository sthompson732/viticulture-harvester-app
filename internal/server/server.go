package server

import (
	"context"
	"viticulture-harvester-app/internal/clients"
	"viticulture-harvester-app/internal/db"
	"viticulture-harvester-app/internal/storage"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server that handles requests.
type Server struct {
	database        *db.DB
	storageClient   *storage.StorageService
	satelliteClient *clients.SatelliteClient
	soilClient      *clients.SoilClient
	router          *gin.Engine
}

// NewServer creates a new instance of Server with initialized routes and dependencies.
func NewServer(database *db.DB, storageClient *storage.StorageService, satelliteClient *clients.SatelliteClient, soilClient *clients.SoilClient) *Server {
	s := &Server{
		database:        database,
		storageClient:   storageClient,
		satelliteClient: satelliteClient,
		soilClient:      soilClient,
		router:          gin.Default(),
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures the API routes. Refactor
func (s *Server) setupRoutes() {
	// Define your routes here
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Example route that uses the database
	s.router.GET("/vineyards", s.handleGetVineyards)

	// Needs refactor
}

// handleGetVineyards is an example handler function that queries the database for vineyards.
func (s *Server) handleGetVineyards(c *gin.Context) {
	// Example: Fetch vineyards from the database
	vineyards, err := s.database.GetVineyards(context.Background()) // Assuming GetVineyards is implemented
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, vineyards)
}

// Start runs the server on a specified port.
func (s *Server) Start(port string) error {
	return s.router.Run(port)
}
