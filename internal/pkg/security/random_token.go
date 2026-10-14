package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRandomToken generates a cryptographically secure random hex string.
// A length of 32 bytes yields a 64-character hex token.
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
