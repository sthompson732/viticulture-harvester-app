/*
 * zoneservice.go: Spatial logic for zone-related operations.
 * Provides functions for generating isochrones and Voronoi polygons.
 * Usage: Handles spatial calculations such as driving zones and influence areas.
 * Author(s): Shannon Thompson
 * Created on: 1/20/2025
 */
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/pzsz/voronoi"
)

type IsochroneResult struct {
	Coordinates [][][]float64 `json:"coordinates"`
}

type ZoneService struct {
	ValhallaURL string // Base URL for Valhalla API
}

func NewZoneService(valhallaURL string) *ZoneService {
	return &ZoneService{
		ValhallaURL: valhallaURL,
	}
}

// GenerateIsochrone fetches driving-time polygons from Valhalla
func (s *ZoneService) GenerateIsochrone(ctx context.Context, lat, lon float64, minutes int) (*IsochroneResult, error) {
	// Construct the Valhalla API URL
	url := fmt.Sprintf("%s/isochrone?json={\"locations\":[{\"lat\":%f,\"lon\":%f}],\"costing\":\"auto\",\"contours\":[{\"time\":%d}]}", s.ValhallaURL, lat, lon, minutes)

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling Valhalla: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Valhalla returned status: %d", resp.StatusCode)
	}

	// Decode the response
	var result struct {
		Features []struct {
			Geometry struct {
				Coordinates [][][]float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Features) > 0 {
		return &IsochroneResult{Coordinates: result.Features[0].Geometry.Coordinates}, nil
	}

	return nil, fmt.Errorf("no isochrone data found")
}

// GenerateVoronoi computes Voronoi polygons for given points
func (s *ZoneService) GenerateVoronoi(points []voronoi.Vertex, bounds orb.Bound) ([]*geojson.Feature, error) {
	diagram := voronoi.ComputeDiagram(points, bounds)
	features := []*geojson.Feature{}

	for _, cell := range diagram.Cells {
		polygon := orb.Polygon{}
		for _, edge := range cell.Halfedges {
			polygon = append(polygon, orb.LineString{
				edge.Start.Point,
				edge.End.Point,
			})
		}
		feature := geojson.NewFeature(polygon)
		features = append(features, feature)
	}
	return features, nil
}
