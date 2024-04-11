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

	"github.com/gorilla/mux"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
	"github.com/sthompson732/viticulture-harvester-app/internal/service"
	"github.com/sthompson732/viticulture-harvester-app/pkg/util"
)

type AppHandler struct {
	VineyardService service.VineyardService
	ImageService    service.ImageService
	SoilDataService service.SoilDataService
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
