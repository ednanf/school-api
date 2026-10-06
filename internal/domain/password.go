package domain

import "errors"

var (
	ErrInvalidPasswordFormat = errors.New("invalid password hash format")
	ErrEmptyPassword         = errors.New("password cannot be empty")
)

// PasswordHasher contain methods to hash and verify passwords
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, encodedHash string) (bool, error)
}
