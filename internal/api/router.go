/*
 * File: router.go
 * Description: Configures and initializes the API routes for the application. This module uses Gorilla Mux
 *              to create a router and define routes that are connected to their respective handlers in the
 *              handlers.go module.
 * Usage:
 *   - Initialize the router and define all route mappings.
 *   - Use middleware for logging and other pre-processing needs.
 * Dependencies:
 *   - Gorilla Mux for route management.
 *   - Handlers for processing requests.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sthompson732/viticulture-harvester-app/internal/service"
)

func NewRouter(vineyardService service.VineyardService, imageService service.ImageService, soilDataService service.SoilDataService) *mux.Router {
	router := mux.NewRouter()

	handler := &AppHandler{
		VineyardService: vineyardService,
		ImageService:    imageService,
		SoilDataService: soilDataService,
	}

	// Middleware to log HTTP requests
	router.Use(loggingMiddleware)

	// Vineyard routes
	router.HandleFunc("/vineyards", handler.CreateVineyard).Methods("POST")
	router.HandleFunc("/vineyards/{id}", handler.GetVineyard).Methods("GET")
	router.HandleFunc("/vineyards/{id}", handler.UpdateVineyard).Methods("PUT")
	router.HandleFunc("/vineyards/{id}", handler.DeleteVineyard).Methods("DELETE")
	router.HandleFunc("/vineyards", handler.ListVineyards).Methods("GET")

	// Image routes
	router.HandleFunc("/images", handler.SaveImage).Methods("POST")
	router.HandleFunc("/images/{id}", handler.GetImage).Methods("GET")
	router.HandleFunc("/images/{id}", handler.DeleteImage).Methods("DELETE")
	router.HandleFunc("/vineyards/{vineyardID}/images", handler.ListImages).Methods("GET")

	// Soil data routes
	router.HandleFunc("/vineyards/{vineyardID}/soil", handler.GetSoilData).Methods("GET")
	router.HandleFunc("/vineyards/{vineyardID}/soil", handler.UpdateSoilData).Methods("PUT")

	return router
}

// Simple middleware for logging HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
