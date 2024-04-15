/*
 * weather.go: Interacts with meteorological data providers to retrieve current and forecasted
 *              weather conditions. This client utilizes APIs to integrate weather data into the
 *              vineyard management system, enhancing decision-making processes with up-to-date
 *              weather information.
 * Usage:
 *   - Retrieve real-time weather data and forecasts.
 *   - Use weather data to adjust vineyard management practices such as irrigation and pest control.
 * Dependencies:
 *   - External weather information services (e.g., OpenWeatherMap, WeatherAPI).
 *   - HTTP client libraries for making API requests.
 * Author(s): Shannon Thompson
 * Created on: 04/14/2024
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

type WeatherClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewWeatherClient(apiKey, baseURL string) *WeatherClient {
	return &WeatherClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (wc *WeatherClient) FetchWeatherData(ctx context.Context, location string) (*model.WeatherData, error) {
	reqURL := fmt.Sprintf("%s/weather?location=%s&apikey=%s", wc.BaseURL, location, wc.APIKey)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := wc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request for weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-ok HTTP status: %s", resp.Status)
	}

	var data model.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &data, nil
}
