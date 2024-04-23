/*
 * router.go: Configures HTTP routes and middleware.
 * Maps endpoints to handler functions, ensuring proper request routing.
 * Usage: Employed by main.go to setup HTTP server routing.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sthompson732/viticulture-harvester-app/internal/config"
	"github.com/sthompson732/viticulture-harvester-app/internal/service"
)

func NewRouter(vineyardService service.VineyardService, imageService service.ImageService,
	soilDataService service.SoilDataService, pestService service.PestService,
	weatherService service.WeatherService, satelliteService service.SatelliteService, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	handler := &AppHandler{
		VineyardService:  vineyardService,
		ImageService:     imageService,
		SoilDataService:  soilDataService,
		PestService:      pestService,
		WeatherService:   weatherService,
		SatelliteService: satelliteService,
	}

	// Middleware for logging and API key verification
	router.Use(loggingMiddleware)
	// Middleware to validate API keys
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if !contains(cfg.ValidAPIKeys, apiKey) {
				http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	defineRoutes(router, handler)

	return router
}

// defineRoutes encapsulates route definitions
func defineRoutes(router *mux.Router, handler *AppHandler) {
	// Vineyard routes
	router.HandleFunc("/vineyards", handler.CreateVineyard).Methods("POST")
	router.HandleFunc("/vineyards/{id}", handler.GetVineyard).Methods("GET")
	router.HandleFunc("/vineyards/{id}", handler.UpdateVineyard).Methods("PUT")
	router.HandleFunc("/vineyards/{id}", handler.DeleteVineyard).Methods("DELETE")
	router.HandleFunc("/vineyards", handler.ListVineyards).Methods("GET")
	router.HandleFunc("/vineyards/{id}/environmental-data", handler.GetVineyardWithEnvironmentalData).Methods("GET")

	// Dynamic route for data fetching based on the data sources defined in config
	router.HandleFunc("/fetch-data", handler.FetchDataFromSource).Methods("GET")

	// Image routes
	router.HandleFunc("/images", handler.SaveImage).Methods("POST")
	router.HandleFunc("/images/{id}", handler.GetImage).Methods("GET")
	router.HandleFunc("/images/{id}", handler.UpdateImage).Methods("PUT")
	router.HandleFunc("/images/{id}", handler.DeleteImage).Methods("DELETE")
	router.HandleFunc("/vineyards/{vineyardID}/images", handler.ListImages).Methods("GET")
	router.HandleFunc("/vineyards/{vineyardID}/images/date-range", handler.FindImagesByDateRange).Methods("POST")
	router.HandleFunc("/vineyards/{vineyardID}/images/recent", handler.GetRecentImages).Methods("GET")

	// Soil data routes
	router.HandleFunc("/soil", handler.CreateSoilData).Methods("POST")
	router.HandleFunc("/soil/{id}", handler.GetSoilData).Methods("GET")
	router.HandleFunc("/soil/{id}", handler.UpdateSoilData).Methods("PUT")
	router.HandleFunc("/soil/{id}", handler.DeleteSoilData).Methods("DELETE")
	router.HandleFunc("/vineyards/{vineyardID}/soil", handler.ListSoilData).Methods("GET")
	router.HandleFunc("/vineyards/{vineyardID}/soil/date-range", handler.ListSoilDataByDateRange).Methods("POST")

	// Pest data routes
	router.HandleFunc("/pests", handler.CreatePestData).Methods("POST")
	router.HandleFunc("/pests/{id}", handler.GetPestData).Methods("GET")
	router.HandleFunc("/pests/{id}", handler.UpdatePestData).Methods("PUT")
	router.HandleFunc("/pests/{id}", handler.DeletePestData).Methods("DELETE")
	router.HandleFunc("/vineyards/{vineyardID}/pests", handler.ListPestData).Methods("GET")
	router.HandleFunc("/vineyards/{vineyardID}/pests/filter", handler.FilterPestData).Methods("POST")

	// Weather data routes
	router.HandleFunc("/weather", handler.CreateWeatherData).Methods("POST")
	router.HandleFunc("/weather/{id}", handler.GetWeatherData).Methods("GET")
	router.HandleFunc("/weather/{id}", handler.UpdateWeatherData).Methods("PUT")
	router.HandleFunc("/weather/{id}", handler.DeleteWeatherData).Methods("DELETE")
	router.HandleFunc("/vineyards/{vineyardID}/weather", handler.ListWeatherData).Methods("GET")
	router.HandleFunc("/vineyards/{vineyardID}/weather/date-range", handler.ListWeatherDataByDateRange).Methods("POST")

	// Satellite routes
	router.HandleFunc("/satellite", handler.CreateSatelliteData).Methods("POST")
	router.HandleFunc("/satellite/{id}", handler.GetSatelliteData).Methods("GET")
	router.HandleFunc("/satellite/{id}", handler.UpdateSatelliteData).Methods("PUT")
	router.HandleFunc("/satellite/{id}", handler.DeleteSatelliteData).Methods("DELETE")
	router.HandleFunc("/vineyards/{vineyardID}/satellite", handler.ListSatelliteData).Methods("GET")
	router.HandleFunc("/vineyards/{vineyardID}/satellite/date-range", handler.ListSatelliteImageryByDateRange).Methods("POST")
	router.HandleFunc("/vineyards/{vineyardID}/satellite/recent", handler.GetRecentSatelliteImages).Methods("GET")
}

// loggingMiddleware logs the HTTP request method and URL path.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Helper function to check if the provided API key is in the list of valid keys
func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
