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

import (
	"time"
)

// Vineyard represents the data model for a vineyard, including its location and soil health.
type Vineyard struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Location         string          `json:"location"`    // Consider using a more complex type for geolocation data
	BoundingBox      string          `json:"boundingBox"` // GeoJSON format for more accurate geospatial representation
	SoilHealth       []SoilData      `json:"soilHealth"`
	SatelliteImagery []SatelliteData `json:"satelliteImagery"`
}

// Image represents metadata about an image related to a vineyard.
type Image struct {
	ID          int       `json:"id"`
	VineyardID  int       `json:"vineyard_id"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CapturedAt  time.Time `json:"capturedAt"`
	BoundingBox string    `json:"boundingBox"` // GeoJSON format to specify the precise area the image covers
}

// SatelliteData represents the structure of data fetched from the satellite imagery API.
type SatelliteData struct {
	ID          int       `json:"id"`
	VineyardID  int       `json:"vineyard_id"`
	ImageURL    string    `json:"imageUrl"`
	CapturedAt  time.Time `json:"capturedAt"`
	Resolution  float64   `json:"resolution"`  // Resolution of the satellite image in meters
	BoundingBox string    `json:"boundingBox"` // GeoJSON format to specify the precise area the satellite image covers
}

// SoilData encapsulates soil characteristics fetched from the soil data API.
type SoilData struct {
	ID               int     `json:"id"`
	VineyardID       int     `json:"vineyard_id"`
	MoistureLevel    float64 `json:"moistureLevel"`
	NutrientContents struct {
		Nitrogen   float64 `json:"nitrogen"`
		Phosphorus float64 `json:"phosphorus"`
		Potassium  float64 `json:"potassium"`
	} `json:"nutrientContents"`
	SoilType  string    `json:"soilType"`
	SampledAt time.Time `json:"sampledAt"`
	Location  string    `json:"location"` // GeoJSON format for precise sampling location
}

// PestData represents data about pest observations within a vineyard.
type PestData struct {
	ID              int       `json:"id"`
	VineyardID      int       `json:"vineyard_id"`
	Description     string    `json:"description"`
	Type            string    `json:"type"`
	Severity        string    `json:"severity"`
	ObservationDate time.Time `json:"observation_date"`
	Location        string    `json:"location"` // GeoJSON format for exact observation location
}

// WeatherData represents weather conditions observed in a vineyard at a specific time.
type WeatherData struct {
	ID              int       `json:"id"`
	VineyardID      int       `json:"vineyard_id"`
	Temperature     float64   `json:"temperature"` // in Celsius
	Humidity        float64   `json:"humidity"`    // percentage
	ObservationTime time.Time `json:"observation_time"`
	Location        string    `json:"location"` // GeoJSON format for exact weather station location
}
