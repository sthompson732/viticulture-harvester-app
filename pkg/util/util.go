/*
 * File: util.go
 * Description: Provides utility functions like JSON response formatting, error handling, and other
 *              common tasks that are used across various parts of the application.
 * Usage:
 *   - Encode objects to JSON for HTTP responses.
 *   - Handle errors consistently across the application.
 *   - Perform common operations like decoding request bodies.
 * Dependencies:
 *   - Standard library packages like 'encoding/json' and 'net/http'.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package util

import (
	"fmt"
	"time"

	"encoding/json"
	"net/http"
)

// JSONResponse simplifies sending JSON responses.
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data) // Ignoring error for brevity; should refactor
	}
}

// ErrorResponse formats and sends an error response.
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]string{"error": message})
}

// DecodeJSONBody is a helper for decoding a JSON request body.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return err
	}
	return nil
}

// ParseDateRange parses start and end date strings and returns time.Time objects.
func ParseDateRange(start, end string) (time.Time, time.Time, error) {
	const layout = "2006-01-02" // ISO 8601 format
	startDate, err := time.Parse(layout, start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date: %w", err)
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("end date %v is before start date %v", endDate, startDate)
	}

	return startDate, endDate, nil
}
