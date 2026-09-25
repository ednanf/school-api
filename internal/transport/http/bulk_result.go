package http

// BulkMeta contains metadata for bulk/batch operations.
type BulkMeta struct {
	Count int `json:"count"` // Number of items processed/created in this batch operation
	Total int `json:"total"` // Overall total, target count, or total requested in the payload
}

// BulkResult represents a generic API response for bulk operations containing items of type T and batch metadata.
type BulkResult[T any] struct {
	Items []T      `json:"items"`
	Meta  BulkMeta `json:"meta"`
}
