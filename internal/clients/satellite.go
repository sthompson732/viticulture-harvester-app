/*
 * File: satellite.go
 * Description: Interacts with satellite data providers to fetch remote sensing and satellite imagery.
 *              Handles API calls and processes responses for use in vineyard monitoring and management.
 * Usage:
 *   - Retrieve satellite images and metadata for specified coordinates and times.
 *   - Process and store imagery data in a compatible format for analysis.
 * Dependencies:
 *   - External satellite imagery APIs (e.g., NASA, ESA).
 *   - JSON parsing and HTTP client libraries.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/config"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

// SatelliteClient is configured to fetch satellite imagery data.
type SatelliteClient struct {
	Config *config.Config
}

// NewSatelliteClient initializes a SatelliteClient with application configuration.
func NewSatelliteClient(cfg *config.Config) *SatelliteClient {
	return &SatelliteClient{
		Config: cfg,
	}
}

// FetchData makes an HTTP request to the satellite imagery API and returns structured data.
func (c *SatelliteClient) FetchData(ctx context.Context, lat, long string, startDate, endDate time.Time) (*model.SatelliteData, error) {
	satelliteConfig, ok := c.Config.DataSources["satellite"]
	if !ok {
		return nil, fmt.Errorf("satellite data source configuration not found")
	}

	reqURL := fmt.Sprintf("%s?lat=%s&lon=%s&start_date=%s&end_date=%s&api_key=%s",
		satelliteConfig.Endpoint,
		lat, long,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		satelliteConfig.APIKey,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var satelliteData model.SatelliteData
	if err := json.NewDecoder(resp.Body).Decode(&satelliteData); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &satelliteData, nil
}
