/*
 * File: handlers.go
 * Description: Provides HTTP handlers for processing API requests. This module contains functions that
 *              interface with the Vineyard, Image, and Soil Data services to handle CRUD operations and
 *              response formatting.
 * Usage:
 *   - Handlers are mapped to specific API endpoints in the router configuration.
 *   - Each handler uses service layer methods to perform business logic.
 * Dependencies:
 *   - Service layer for business logic.
 *   - Gorilla Mux for routing.
 *   - Utility functions for JSON encoding and decoding.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
	"github.com/sthompson732/viticulture-harvester-app/internal/service"
	"github.com/sthompson732/viticulture-harvester-app/pkg/util"
)

type AppHandler struct {
	VineyardService  service.VineyardService
	ImageService     service.ImageService
	SoilDataService  service.SoilDataService
	PestService      service.PestService
	WeatherService   service.WeatherService
	SatelliteService service.SatelliteService
}

// CreateVineyard handles POST requests to add new vineyards
func (h *AppHandler) CreateVineyard(w http.ResponseWriter, r *http.Request) {
	var vineyard model.Vineyard
	if err := json.NewDecoder(r.Body).Decode(&vineyard); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}
	err := h.VineyardService.CreateVineyard(r.Context(), &vineyard)
	if err != nil {
		http.Error(w, "Failed to create vineyard", 500)
		return
	}
	json.NewEncoder(w).Encode(vineyard)
}

// GetVineyard handles GET requests for retrieving a single vineyard by ID
func (h *AppHandler) GetVineyard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	vineyard, err := h.VineyardService.GetVineyard(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch vineyard", 500)
		return
	}
	json.NewEncoder(w).Encode(vineyard)
}

// UpdateVineyard handles PUT requests to update a vineyard by ID
func (h *AppHandler) UpdateVineyard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var vineyard model.Vineyard
	if err := json.NewDecoder(r.Body).Decode(&vineyard); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}
	err := h.VineyardService.UpdateVineyard(r.Context(), id, &vineyard)
	if err != nil {
		http.Error(w, "Failed to update vineyard", 500)
		return
	}
	json.NewEncoder(w).Encode(vineyard)
}

// DeleteVineyard handles DELETE requests to remove a vineyard by ID
func (h *AppHandler) DeleteVineyard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.VineyardService.DeleteVineyard(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete vineyard", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ListVineyards handles GET requests to list all vineyards
func (h *AppHandler) ListVineyards(w http.ResponseWriter, r *http.Request) {
	vineyards, err := h.VineyardService.ListVineyards(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch vineyards", 500)
		return
	}
	json.NewEncoder(w).Encode(vineyards)
}

// GetVineyardWithEnvironmentalData retrieves a vineyard along with its related satellite imagery and soil data.
func (h *AppHandler) GetVineyardWithEnvironmentalData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	vineyard, err := h.VineyardService.GetVineyardWithEnvironmentalData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyard)
}

// Handlers for Images

func (h *AppHandler) SaveImage(w http.ResponseWriter, r *http.Request) {
	var image model.SatelliteData
	if err := util.DecodeJSONBody(w, r, &image); err != nil {
		return
	}
	err := h.ImageService.SaveImage(r.Context(), &image)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not save image")
		return
	}
	util.JSONResponse(w, http.StatusCreated, image)
}

func (h *AppHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	image, err := h.ImageService.GetImage(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch image")
		return
	}
	util.JSONResponse(w, http.StatusOK, image)
}

func (h *AppHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := h.ImageService.DeleteImage(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not delete image")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *AppHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	vineyardID := mux.Vars(r)["vineyardID"]
	images, err := h.ImageService.ListImages(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list images")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}

// UpdateImage updates an existing image record.
func (h *AppHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	var image model.Image
	if err := util.DecodeJSONBody(w, r, &image); err != nil {
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid image ID")
		return
	}
	image.ID = id

	if err := h.ImageService.UpdateImage(r.Context(), &image); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update image")
		return
	}
	util.JSONResponse(w, http.StatusOK, nil)
}

// FindImagesByDateRange retrieves images for a vineyard within a specified date range.
func (h *AppHandler) FindImagesByDateRange(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	// Assuming dates are provided as query parameters
	start, end, err := util.ParseDateRange(r.URL.Query().Get("start"), r.URL.Query().Get("end"))
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid date range")
		return
	}

	images, err := h.ImageService.FindImagesByDateRange(r.Context(), vineyardID, start, end)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not find images")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}

// GetRecentImages retrieves the most recent images for a vineyard.
func (h *AppHandler) GetRecentImages(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5 // default limit to 5 images if not specified or on error
	}

	images, err := h.ImageService.GetRecentImages(r.Context(), vineyardID, limit)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not get recent images")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}

// Handlers for Soil Data

func (h *AppHandler) GetSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID := mux.Vars(r)["vineyardID"]
	soilData, err := h.SoilDataService.GetSoilData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, soilData)
}

func (h *AppHandler) UpdateSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID := mux.Vars(r)["vineyardID"]
	var soilData model.SoilData
	if err := util.DecodeJSONBody(w, r, &soilData); err != nil {
		return
	}
	err := h.SoilDataService.UpdateSoilData(r.Context(), vineyardID, &soilData)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "soil data updated"})
}

// CreateSoilData handles the creation of a new soil data record.
func (h *AppHandler) CreateSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	var soilData model.SoilData
	if err := util.DecodeJSONBody(w, r, &soilData); err != nil {
		return
	}
	soilData.VineyardID = vineyardID

	if err := h.SoilDataService.CreateSoilData(r.Context(), &soilData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not create soil data")
		return
	}
	util.JSONResponse(w, http.StatusCreated, soilData)
}

// DeleteSoilData handles the deletion of a soil data record.
func (h *AppHandler) DeleteSoilData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid soil data ID")
		return
	}

	if err := h.SoilDataService.DeleteSoilData(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not delete soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, nil)
}

// ListSoilData retrieves all soil data entries for a specified vineyard.
func (h *AppHandler) ListSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	soils, err := h.SoilDataService.ListSoilData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, soils)
}

// ListSoilDataByDateRange retrieves soil data within a specified date range for a vineyard.
func (h *AppHandler) ListSoilDataByDateRange(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	start, end, err := util.ParseDateRange(r.URL.Query().Get("start"), r.URL.Query().Get("end"))
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid date range")
		return
	}

	soils, err := h.SoilDataService.ListSoilDataByDateRange(r.Context(), vineyardID, start, end)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not retrieve soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, soils)
}

// Pest Handlers

func (h *AppHandler) GetPestData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid pest ID", http.StatusBadRequest)
		return
	}
	data, err := h.PestService.GetPestData(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch pest data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *AppHandler) CreatePestData(w http.ResponseWriter, r *http.Request) {
	var data model.PestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.PestService.CreatePestData(r.Context(), &data); err != nil {
		http.Error(w, "Failed to create pest data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *AppHandler) DeletePestData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid pest ID", http.StatusBadRequest)
		return
	}
	if err := h.PestService.DeletePestData(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete pest data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AppHandler) ListPests(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		http.Error(w, "Invalid vineyard ID", http.StatusBadRequest)
		return
	}
	pests, err := h.PestService.ListPests(r.Context(), vineyardID)
	if err != nil {
		http.Error(w, "Failed to list pests", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pests)
}

// UpdatePestData handles the updating of pest data records.
func (h *AppHandler) UpdatePestData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid pest data ID")
		return
	}

	var pestData model.PestData
	if err := util.DecodeJSONBody(w, r, &pestData); err != nil {
		return
	}
	pestData.ID = id

	if err := h.PestService.UpdatePestData(r.Context(), &pestData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update pest data")
		return
	}
	util.JSONResponse(w, http.StatusOK, nil)
}

// ListPestData retrieves all pest data entries for a specified vineyard.
func (h *AppHandler) ListPestData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	pests, err := h.PestService.ListPestData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list pest data")
		return
	}
	util.JSONResponse(w, http.StatusOK, pests)
}

// FilterPestData retrieves pest data based on type and severity filters.
func (h *AppHandler) FilterPestData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	pestType := r.URL.Query().Get("type")
	severity := r.URL.Query().Get("severity")

	pests, err := h.PestService.FilterPestData(r.Context(), vineyardID, pestType, severity)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not filter pest data")
		return
	}
	util.JSONResponse(w, http.StatusOK, pests)
}

// UpdateWeatherData handles the updating of weather data records.
func (h *AppHandler) UpdateWeatherData(w http.ResponseWriter, r *http.Request) {
	var weatherData model.WeatherData
	if err := util.DecodeJSONBody(w, r, &weatherData); err != nil {
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid weather data ID")
		return
	}
	weatherData.ID = id

	if err := h.WeatherService.UpdateWeatherData(r.Context(), &weatherData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, nil)
}

// ListWeatherDataByDateRange retrieves weather data within a specified date range for a vineyard.
func (h *AppHandler) ListWeatherDataByDateRange(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	start, end, err := util.ParseDateRange(r.URL.Query().Get("start"), r.URL.Query().Get("end"))
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid date range")
		return
	}

	weather, err := h.WeatherService.ListWeatherDataByDateRange(r.Context(), vineyardID, start, end)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not retrieve weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, weather)
}

// CreateSatelliteData handles the creation of new satellite imagery data records.
func (h *AppHandler) CreateSatelliteData(w http.ResponseWriter, r *http.Request) {
	var satelliteData model.SatelliteData
	if err := util.DecodeJSONBody(w, r, &satelliteData); err != nil {
		return
	}
	if err := h.SatelliteService.CreateSatelliteData(r.Context(), &satelliteData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not create satellite data")
		return
	}
	util.JSONResponse(w, http.StatusCreated, satelliteData)
}

// GetSatelliteData retrieves a single satellite data record by ID.
func (h *AppHandler) GetSatelliteData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid satellite data ID")
		return
	}

	satelliteData, err := h.SatelliteService.GetSatelliteData(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, satelliteData)
}

// Weather Handlers

func (h *AppHandler) GetWeatherData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid weather ID", http.StatusBadRequest)
		return
	}
	data, err := h.WeatherService.GetWeatherData(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *AppHandler) CreateWeatherData(w http.ResponseWriter, r *http.Request) {
	var data model.WeatherData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.WeatherService.CreateWeatherData(r.Context(), &data); err != nil {
		http.Error(w, "Failed to create weather data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *AppHandler) DeleteWeatherData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid weather ID", http.StatusBadRequest)
		return
	}
	if err := h.WeatherService.DeleteWeatherData(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete weather data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AppHandler) ListWeatherData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		http.Error(w, "Invalid vineyard ID", http.StatusBadRequest)
		return
	}
	weatherData, err := h.WeatherService.ListWeatherData(r.Context(), vineyardID)
	if err != nil {
		http.Error(w, "Failed to list weather data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(weatherData)
}

// UpdateSatelliteData handles the updating of satellite data records.
func (h *AppHandler) UpdateSatelliteData(w http.ResponseWriter, r *http.Request) {
	var satelliteData model.SatelliteData
	if err := util.DecodeJSONBody(w, r, &satelliteData); err != nil {
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid satellite data ID")
		return
	}
	satelliteData.ID = id

	if err := h.SatelliteService.UpdateSatelliteData(r.Context(), &satelliteData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, nil)
}

// DeleteSatelliteData handles the deletion of a satellite data record.
func (h *AppHandler) DeleteSatelliteData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid satellite data ID")
		return
	}

	if err := h.SatelliteService.DeleteSatelliteData(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not delete satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, nil)
}

// ListSatelliteData retrieves all satellite data entries for a specified vineyard.
func (h *AppHandler) ListSatelliteData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	data, err := h.SatelliteService.ListSatelliteData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, data)
}

// ListSatelliteImageryByDateRange retrieves satellite imagery within a specified date range for a vineyard.
func (h *AppHandler) ListSatelliteImageryByDateRange(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	start, end, err := util.ParseDateRange(r.URL.Query().Get("start"), r.URL.Query().Get("end"))
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid date range")
		return
	}

	images, err := h.SatelliteService.ListSatelliteImageryByDateRange(r.Context(), vineyardID, start, end)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not retrieve satellite imagery")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}

// GetRecentSatelliteImages retrieves the most recent satellite images for a vineyard.
func (h *AppHandler) GetRecentSatelliteImages(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5 // default limit to 5 images if not specified or on error
	}

	images, err := h.SatelliteService.GetRecentSatelliteImages(r.Context(), vineyardID, limit)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not get recent satellite images")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}
