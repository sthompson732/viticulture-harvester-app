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
