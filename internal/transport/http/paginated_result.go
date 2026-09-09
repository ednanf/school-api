package http

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
