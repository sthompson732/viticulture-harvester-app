/*
 * File: models.go
 * Description: Defines data structures that represent database tables and are used throughout the application
 *              for data manipulation and retrieval. These models are used directly with the database/sql
 *              package to prepare and execute SQL statements.
 * Usage:
 *   - Structs are used to scan results from SQL queries and to structure data for insertion.
 *   - Serve as a data transfer object between the database and application logic.
 * Dependencies:
 *   - Used directly by db.go for constructing SQL queries and scanning query results.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package model

import "time"

// SatelliteData represents the structure of data fetched from the satellite imagery API.
type SatelliteData struct {
	ImageURL   string    `json:"imageUrl"`
	CapturedAt time.Time `json:"capturedAt"`
	// Add more fields as per the satellite API response structure
}

// SoilData encapsulates soil characteristics fetched from the soil data API.
type SoilData struct {
	MoistureLevel    float64 `json:"moistureLevel"`
	NutrientContents struct {
		Nitrogen   float64 `json:"nitrogen"`
		Phosphorus float64 `json:"phosphorus"`
		Potassium  float64 `json:"potassium"`
	} `json:"nutrientContents"`
	SoilType string `json:"soilType"`
	// Additional soil properties as needed
}

// Vineyard represents the data model for a vineyard, including its location and soil health.
type Vineyard struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Location         string `json:"location"` // Consider using a more complex type for geolocation data
	SoilHealth       []SoilData
	SatelliteImagery []SatelliteData
	// Consider relationships, such as vines or crops, that we may also want to be part of this model
}

// PestData represents data about pest observations within a vineyard.
type PestData struct {
	ID              int       `json:"id"`
	VineyardID      int       `json:"vineyard_id"`
	Description     string    `json:"description"`
	ObservationDate time.Time `json:"observation_date"`
	Location        string    `json:"location"` // GeoJSON format
}

// WeatherData represents weather conditions observed in a vineyard at a specific time.
type WeatherData struct {
	ID              int       `json:"id"`
	VineyardID      int       `json:"vineyard_id"`
	Temperature     float64   `json:"temperature"` // in Celsius
	Humidity        float64   `json:"humidity"`    // percentage
	ObservationTime time.Time `json:"observation_time"`
	Location        string    `json:"location"` // GeoJSON format
}
