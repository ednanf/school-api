package http

// BatchMeta contains metadata for bulk/batch operations.
type BatchMeta struct {
	Count int `json:"count"` // Number of items processed/created in this batch operation
	Total int `json:"total"` // Overall total, target count, or total requested in the payload
}

// BatchResult represents a generic API response for bulk operations containing items of type T and batch metadata.
type BatchResult[T any] struct {
	Items []T       `json:"items"`
	Meta  BatchMeta `json:"meta"`
}
