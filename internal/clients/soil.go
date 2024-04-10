package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"viticulture-harvester-app/internal/config"
	"viticulture-harvester-app/internal/model"
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
	// Constructing the request URL with parameters from the configuration and function arguments
	reqURL := fmt.Sprintf("%s?latitude=%s&longitude=%s&api_key=%s",
		c.Config.DataSources.Soil.APIEndpoint,
		lat, long,
		c.Config.DataSources.Soil.APIKey,
	)

	// Making the HTTP request with context for cancellation and timeouts
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

	// Verifying response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Reading and decoding the response body into the SoilData struct
	var soilData model.SoilData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	if err := json.Unmarshal(body, &soilData); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &soilData, nil
}
