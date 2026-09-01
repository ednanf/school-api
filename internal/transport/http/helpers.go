package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// sendJSON is a quick helper to reduce boilerplate, acting like res.json()
func sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// formatValidationErrors is a quick helper to map validation errors into a clean key-value format for response bodies
func formatValidationErrors(errs validator.ValidationErrors) map[string]string {
	details := make(map[string]string)
	for _, err := range errs {
		// Example output: {"Email": "email", "FirstName": "required"}
		details[err.Field()] = err.Tag()
	}
	return details
}
