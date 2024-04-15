/*
 * models.go: Defines data structures for the application.
 * Structures are used for scanning SQL results and preparing data for transactions.
 * Usage: Serves as a transfer object between the database and application logic.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package model

import (
	"io"
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
	FilePath    string    `json:"filePath"`    // Local or remote file path of the image for uploading
	ImageFile   io.Reader `json:"-"`           // The image file data, excluded from JSON operations
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
	Location  Location  `json:"location"` // Modified to use a structured type
}

// Location struct to hold geospatial coordinates
type Location struct {
	X float64 `json:"longitude"` // longitude
	Y float64 `json:"latitude"`  // latitude
}

// PestData represents data about pest observations within a vineyard.
type PestData struct {
	ID              int       `json:"id"`
	VineyardID      int       `json:"vineyard_id"`
	Description     string    `json:"description"`
	Type            string    `json:"type"`
	Severity        string    `json:"severity"`
	ObservationDate time.Time `json:"observation_date"`
	Location        Location  `json:"location"` // Modified to use a structured type
}

// WeatherData represents weather conditions observed in a vineyard at a specific time.
type WeatherData struct {
	ID              int       `json:"id"`
	VineyardID      int       `json:"vineyard_id"`
	Temperature     float64   `json:"temperature"` // in Celsius
	Humidity        float64   `json:"humidity"`    // percentage
	ObservationTime time.Time `json:"observation_time"`
	Location        Location  `json:"location"` // Modified to use a structured type
}
