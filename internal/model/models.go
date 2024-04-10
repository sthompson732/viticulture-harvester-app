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
	SoilHealth       SoilData
	SatelliteImagery []SatelliteData
	// Consider relationships, such as vines or crops, that may also be part of this model
}
