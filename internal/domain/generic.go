package domain

type PaginatedResult[T any] struct {
	Total int `json:"total"`
	Items []T `json:"items"`
}
