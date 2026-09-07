package domain

// PaginatedResult represents a paginated API response containing a slice of items of type T
// and the total count of available items matching the request.
type PaginatedResult[T any] struct {
	// Total is the total number of items available across all pages.
	Total int `json:"total"`

	// Items is the list of records returned for the current page.
	Items []T `json:"items"`
}
