/*
 * File: pest.go
 * Description: Manages interactions with pest control data providers to fetch pest activity data
 *              and pest management solutions. This client handles all API calls to pest control services,
 *              consolidating pest data retrieval and analysis for use in vineyard management.
 * Usage:
 *   - Fetch pest activity reports and recommendations for specific areas.
 *   - Analyze pest data to provide actionable insights and management solutions.
 * Dependencies:
 *   - External pest control data APIs (e.g., local agricultural extensions or pest control technology providers).
 *   - JSON parsing and HTTP client libraries for making API requests.
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

	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type PestClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewPestClient(apiKey, baseURL string) *PestClient {
	return &PestClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (pc *PestClient) FetchPestData(ctx context.Context, location string) (*model.PestData, error) {
	reqURL := fmt.Sprintf("%s/api/pests?location=%s&apikey=%s", pc.BaseURL, location, pc.APIKey)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request for pest data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-ok HTTP status: %s", resp.Status)
	}

	var data model.PestData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &data, nil
}
