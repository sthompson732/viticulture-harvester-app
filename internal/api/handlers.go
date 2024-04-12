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

// Handlers for Vineyard
// CreateVineyard handles POST requests to add new vineyards

func (h *AppHandler) CreateVineyard(w http.ResponseWriter, r *http.Request) {
	var vineyard model.Vineyard
	if err := json.NewDecoder(r.Body).Decode(&vineyard); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.VineyardService.CreateVineyard(r.Context(), &vineyard); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to create vineyard")
		return
	}
	util.JSONResponse(w, http.StatusCreated, vineyard)
}

// GetVineyard handles GET requests for retrieving a single vineyard by ID
func (h *AppHandler) GetVineyard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	vineyard, err := h.VineyardService.GetVineyard(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyard)
}

// UpdateVineyard handles PUT requests to update a vineyard by ID
func (h *AppHandler) UpdateVineyard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	var vineyard model.Vineyard
	if err := json.NewDecoder(r.Body).Decode(&vineyard); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.VineyardService.UpdateVineyard(r.Context(), id, &vineyard); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to update vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyard)
}

// DeleteVineyard handles DELETE requests to remove a vineyard by ID
func (h *AppHandler) DeleteVineyard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	if err := h.VineyardService.DeleteVineyard(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ListVineyards handles GET requests to list all vineyards
func (h *AppHandler) ListVineyards(w http.ResponseWriter, r *http.Request) {
	vineyards, err := h.VineyardService.ListVineyards(r.Context())
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch vineyards")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyards)
}

// GetVineyardWithEnvironmentalData retrieves a vineyard along with its related satellite imagery and soil data.
func (h *AppHandler) GetVineyardWithEnvironmentalData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	vineyard, err := h.VineyardService.GetVineyardWithEnvironmentalData(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch vineyard with environmental data")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyard)
}

// Handlers for Images

func (h *AppHandler) SaveImage(w http.ResponseWriter, r *http.Request) {
	var image model.Image
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.ImageService.SaveImage(r.Context(), &image); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not save image")
		return
	}
	util.JSONResponse(w, http.StatusCreated, image)
}

func (h *AppHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid image ID")
		return
	}
	image, err := h.ImageService.GetImage(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch image")
		return
	}
	util.JSONResponse(w, http.StatusOK, image)
}

func (h *AppHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid image ID")
		return
	}
	var image model.Image
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	image.ID = id
	if err := h.ImageService.UpdateImage(r.Context(), &image); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update image")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *AppHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid image ID")
		return
	}
	if err := h.ImageService.DeleteImage(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not delete image")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
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
		limit = 5 // Default to 5 images if limit is not specified or on error
	}
	images, err := h.ImageService.GetRecentImages(r.Context(), vineyardID, limit)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not get recent images")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}

func (h *AppHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	images, err := h.ImageService.ListImages(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list images")
		return
	}

	util.JSONResponse(w, http.StatusOK, images)
}

// Handlers for Soil Data
// CreateSoilData handles the creation of a new soil data record.

func (h *AppHandler) CreateSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	var soilData model.SoilData
	if err := json.NewDecoder(r.Body).Decode(&soilData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	soilData.VineyardID = vineyardID

	if err := h.SoilDataService.CreateSoilData(r.Context(), &soilData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not create soil data")
		return
	}
	util.JSONResponse(w, http.StatusCreated, soilData)
}

func (h *AppHandler) GetSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	soilData, err := h.SoilDataService.GetSoilData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, soilData)
}

func (h *AppHandler) UpdateSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	var soilData model.SoilData
	if err := json.NewDecoder(r.Body).Decode(&soilData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	soilData.VineyardID = vineyardID

	if err := h.SoilDataService.UpdateSoilData(r.Context(), vineyardID, &soilData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, soilData)
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
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ListSoilData retrieves all soil data entries for a specified vineyard.
func (h *AppHandler) ListSoilData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	soilData, err := h.SoilDataService.ListSoilData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list soil data")
		return
	}
	util.JSONResponse(w, http.StatusOK, soilData)
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

// Handlers for Pest Data

func (h *AppHandler) CreatePestData(w http.ResponseWriter, r *http.Request) {
	var pestData model.PestData
	if err := json.NewDecoder(r.Body).Decode(&pestData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.PestService.CreatePestData(r.Context(), &pestData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to create pest data")
		return
	}
	util.JSONResponse(w, http.StatusCreated, pestData)
}

func (h *AppHandler) GetPestData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid pest data ID")
		return
	}
	pestData, err := h.PestService.GetPestData(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch pest data")
		return
	}
	util.JSONResponse(w, http.StatusOK, pestData)
}

// UpdatePestData handles the updating of pest data records.

func (h *AppHandler) UpdatePestData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid pest data ID")
		return
	}
	var pestData model.PestData
	if err := json.NewDecoder(r.Body).Decode(&pestData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	pestData.ID = id
	if err := h.PestService.UpdatePestData(r.Context(), &pestData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update pest data")
		return
	}
	util.JSONResponse(w, http.StatusOK, pestData)
}

func (h *AppHandler) DeletePestData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid pest data ID")
		return
	}
	if err := h.PestService.DeletePestData(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete pest data")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *AppHandler) ListPests(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	pests, err := h.PestService.ListPests(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to list pests")
		return
	}
	util.JSONResponse(w, http.StatusOK, pests)
}

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

// Handlers for Weather Data

func (h *AppHandler) CreateWeatherData(w http.ResponseWriter, r *http.Request) {
	var weatherData model.WeatherData
	if err := json.NewDecoder(r.Body).Decode(&weatherData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.WeatherService.CreateWeatherData(r.Context(), &weatherData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to create weather data")
		return
	}
	util.JSONResponse(w, http.StatusCreated, weatherData)
}

func (h *AppHandler) GetWeatherData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid weather data ID")
		return
	}
	weatherData, err := h.WeatherService.GetWeatherData(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, weatherData)
}

func (h *AppHandler) UpdateWeatherData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid weather data ID")
		return
	}
	var weatherData model.WeatherData
	if err := json.NewDecoder(r.Body).Decode(&weatherData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	weatherData.ID = id
	if err := h.WeatherService.UpdateWeatherData(r.Context(), &weatherData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *AppHandler) DeleteWeatherData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid weather data ID")
		return
	}
	if err := h.WeatherService.DeleteWeatherData(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *AppHandler) ListWeatherData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}

	weatherData, err := h.WeatherService.ListWeatherDataByVineyard(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, weatherData)
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
	weatherData, err := h.WeatherService.ListWeatherDataByDateRange(r.Context(), vineyardID, start, end)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not retrieve weather data")
		return
	}
	util.JSONResponse(w, http.StatusOK, weatherData)
}

// Handlers for Satellite Data
// CreateSatelliteData handles the creation of new satellite imagery data records.
func (h *AppHandler) CreateSatelliteData(w http.ResponseWriter, r *http.Request) {
	var satelliteData model.SatelliteData
	if err := json.NewDecoder(r.Body).Decode(&satelliteData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.SatelliteService.CreateSatelliteData(r.Context(), &satelliteData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to create satellite data")
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
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, satelliteData)
}

// UpdateSatelliteData handles the updating of satellite data records.
func (h *AppHandler) UpdateSatelliteData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid satellite data ID")
		return
	}
	var satelliteData model.SatelliteData
	if err := json.NewDecoder(r.Body).Decode(&satelliteData); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	satelliteData.ID = id
	if err := h.SatelliteService.UpdateSatelliteData(r.Context(), &satelliteData); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteSatelliteData handles the deletion of a satellite data record.
func (h *AppHandler) DeleteSatelliteData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid satellite data ID")
		return
	}
	if err := h.SatelliteService.DeleteSatelliteData(r.Context(), id); err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete satellite data")
		return
	}

	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ListSatelliteData retrieves all satellite data entries for a specified vineyard.
func (h *AppHandler) ListSatelliteData(w http.ResponseWriter, r *http.Request) {
	vineyardID, err := strconv.Atoi(mux.Vars(r)["vineyardID"])
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid vineyard ID")
		return
	}
	satelliteData, err := h.SatelliteService.ListSatelliteData(r.Context(), vineyardID)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list satellite data")
		return
	}
	util.JSONResponse(w, http.StatusOK, satelliteData)
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
	imagery, err := h.SatelliteService.ListSatelliteImageryByDateRange(r.Context(), vineyardID, start, end)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not retrieve satellite imagery")
		return
	}
	util.JSONResponse(w, http.StatusOK, imagery)
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
		limit = 5 // Default to 5 images if limit is not specified or on error
	}
	images, err := h.SatelliteService.GetRecentSatelliteImages(r.Context(), vineyardID, limit)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not get recent satellite images")
		return
	}
	util.JSONResponse(w, http.StatusOK, images)
}
