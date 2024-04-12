/*
 * File: soil.go
 * Description: Retrieves soil health data from various soil data services, enabling detailed soil analysis
 *              for precision agriculture applications in vineyards.
 * Usage:
 *   - Query soil composition, moisture levels, and other relevant soil health indicators.
 *   - Use data to inform soil management practices and interventions.
 * Dependencies:
 *   - Soil data providers (e.g., SoilGrids, local agricultural services).
 *   - Libraries for handling HTTP requests and JSON data.
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

// SoilClient manages interactions with a soil data API.
type SoilClient struct {
	Config *config.Config
}

// NewSoilClient creates a new instance of SoilClient with configuration settings.
func NewSoilClient(cfg *config.Config) *SoilClient {
	return &SoilClient{
		Config: cfg,
	}
}

// FetchData queries the soil data API and returns structured information.
func (c *SoilClient) FetchData(ctx context.Context, lat, long string) (*model.SoilData, error) {
	soilConfig, ok := c.Config.DataSources["soil"]
	if !ok {
		return nil, fmt.Errorf("soil data source configuration not found")
	}

	reqURL := fmt.Sprintf("%s?latitude=%s&longitude=%s&api_key=%s",
		soilConfig.Endpoint,
		lat, long,
		soilConfig.APIKey,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
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

	var soilData model.SoilData
	if err := json.NewDecoder(resp.Body).Decode(&soilData); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &soilData, nil
}
