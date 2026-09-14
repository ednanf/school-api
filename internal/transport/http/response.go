package http

import (
	"encoding/json"
	"net/http"
)

// PaginatedMeta contains metadata about the paginated dataset.
type PaginatedMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Count      int `json:"count"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResult represents a generic paginated API response wrapping items of type T and response metadata.
type PaginatedResult[T any] struct {
	Items []T           `json:"items"`
	Meta  PaginatedMeta `json:"meta"`
}

// APIResponse defines the standard JSON envelope
type APIResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

// SuccessPayload enforces that every success response has a message and data
type SuccessPayload struct {
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"` // omitempty hides field if nil (e.g. for simple 200 OK)
}

// ErrorPayload enforces message + optional error details
type ErrorPayload struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// sendSuccess handles all 2xx responses
func sendSuccess(w http.ResponseWriter, statusCode int, message string, result any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Data: SuccessPayload{
			Message: message,
			Result:  result,
		},
	})
}

// sendError handles 4xx and 5xx error responses
func sendError(w http.ResponseWriter, statusCode int, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Status: "error",
		Data: ErrorPayload{
			Message: message,
			Details: details,
		},
	})
}
