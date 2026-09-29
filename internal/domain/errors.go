package domain

import "errors"

var (
	// Important to decouple SQL from layers - eliminating the necessity of importing for sql.ErrNoRows
	ErrNotFound = errors.New("resource not found")
)
