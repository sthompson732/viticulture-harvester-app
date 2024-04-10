package api

import (
	"net/http"
	"viticulture-harvester-app/internal/model"
	"viticulture-harvester-app/internal/service"
	"viticulture-harvester-app/internal/util"

	"github.com/gorilla/mux"
)

// AppHandler struct encapsulates dependencies for handlers
type AppHandler struct {
	VineyardService service.VineyardService
	ImageService    service.ImageService
	SoilDataService service.SoilDataService
}

// Handler for creating a vineyard
func (h *AppHandler) CreateVineyard(w http.ResponseWriter, r *http.Request) {
	var vineyard model.Vineyard
	if err := util.DecodeJSONBody(w, r, &vineyard); err != nil {
		return
	}
	err := h.VineyardService.CreateVineyard(r.Context(), &vineyard)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not save vineyard")
		return
	}
	util.JSONResponse(w, http.StatusCreated, vineyard)
}

// Handler for retrieving a single vineyard by ID
func (h *AppHandler) GetVineyard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	vineyard, err := h.VineyardService.GetVineyard(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyard)
}

// Handler for updating a vineyard by ID
func (h *AppHandler) UpdateVineyard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var vineyard model.Vineyard
	if err := util.DecodeJSONBody(w, r, &vineyard); err != nil {
		return
	}
	err := h.VineyardService.UpdateVineyard(r.Context(), id, &vineyard)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not update vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyard)
}

// Handler for deleting a vineyard by ID
func (h *AppHandler) DeleteVineyard(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := h.VineyardService.DeleteVineyard(r.Context(), id)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not delete vineyard")
		return
	}
	util.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// Handler for listing all vineyards
func (h *AppHandler) ListVineyards(w http.ResponseWriter, r *http.Request) {
	vineyards, err := h.VineyardService.ListVineyards(r.Context())
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Could not list vineyards")
		return
	}
	util.JSONResponse(w, http.StatusOK, vineyards)
}

// Handlers for Images

func (h *AppHandler) SaveImage(w http.ResponseWriter, r *http.Request) {
	var image model.Image
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
