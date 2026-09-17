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

// APIResponse defines the standard JSON envelope with optional fields.
type APIResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message,omitempty"` // Omitted for standard 200 OK
	Data    any            `json:"data,omitempty"`
	Meta    *PaginatedMeta `json:"meta,omitempty"`
	Details any            `json:"details,omitempty"`
}

// sendSuccess handles standard success responses.
// Automatically omits "message" for 200 OK responses.
func sendSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	resp := APIResponse{
		Status: "success",
		Data:   data,
	}

	// Only include message for non-200 OK success statuses (e.g., 201 Created)
	if statusCode != http.StatusOK {
		resp.Message = message
	}

	writeJSON(w, statusCode, resp)
}

// sendPaginated handles paginated array responses (typically 200 OK).
func sendPaginated[T any](w http.ResponseWriter, statusCode int, items []T, meta PaginatedMeta) {
	writeJSON(w, statusCode, APIResponse{
		Status: "success",
		Data:   items,
		Meta:   &meta,
	})
}

// sendError handles 4xx and 5xx error responses.
func sendError(w http.ResponseWriter, statusCode int, message string, details any) {
	writeJSON(w, statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Details: details,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
