package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"viticulture-harvester-app/internal/config"
	"viticulture-harvester-app/internal/model"
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
	// Constructing the request URL from the configuration and method parameters
	reqURL := fmt.Sprintf("%s?lat=%s&lon=%s&start_date=%s&end_date=%s&api_key=%s",
		c.Config.DataSources.Satellite.APIEndpoint,
		lat, long,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		c.Config.DataSources.Satellite.APIKey,
	)

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	// Execute the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the JSON response body into the SatelliteData struct
	var satelliteData model.SatelliteData
	if err := json.NewDecoder(resp.Body).Decode(&satelliteData); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &satelliteData, nil
}
