package domain

import "errors"

var (
	// Important to decouple SQL from layers - eliminating the necessity of importing for sql.ErrNoRows
	ErrNotFound = errors.New("resource not found")

	ErrUserNotFound = errors.New("user not found")
	ErrUserInactive = errors.New("user account is deactivated")

	ErrEmailAlreadyExists    = errors.New("email already in use")
	ErrUsernameAlreadyExists = errors.New("username already in use")

	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
