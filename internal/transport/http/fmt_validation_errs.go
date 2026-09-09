package http

import (
	"github.com/go-playground/validator/v10"
)

// formatValidationErrors is a quick helper to map validation errors into a clean key-value format for response bodies
func formatValidationErrors(errs validator.ValidationErrors) map[string]string {
	details := make(map[string]string)
	for _, err := range errs {
		// Example output: {"Email": "email", "FirstName": "required"}
		details[err.Field()] = err.Tag()
	}
	return details
}
