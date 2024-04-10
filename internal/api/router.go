package api

import (
	"log"
	"net/http"
	"viticulture-harvester-app/internal/service"

	"github.com/gorilla/mux"
)

// NewRouter initializes a new router with all the route mappings
func NewRouter(vineyardService service.VineyardService, imageService service.ImageService, soilDataService service.SoilDataService) *mux.Router {
	handler := &AppHandler{
		VineyardService: vineyardService,
		ImageService:    imageService,
		SoilDataService: soilDataService,
	}

	router := mux.NewRouter()

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

	// Middleware to log HTTP requests
	router.Use(loggingMiddleware)

	return router
}

// Simple middleware for logging HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request, e.g., method and URL
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
