package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/crypto/argon2"
)

type argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

// NewArgon2Hasher initializes sensible, OWASP-aligned defaults for web apps:
// 64 MB memory, 1 iteration, 4 parallel threads, 16-byte salt, 32-byte key.
func NewArgon2Hasher() domain.PasswordHasher {
	// By returning the interface domain.PasswordHasher, the contract is enforced right at the creation
	return &argon2Hasher{
		time:    1,
		memory:  64 * 1024, // 64 MB
		threads: 4,
		keyLen:  32,
		saltLen: 16,
	}
}

// Hash receives a password string and returns a hashed string with 16-byte salt and 32-byte key
// encoded to a standard PHC format:
// `$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>`
func (h *argon2Hasher) Hash(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", domain.ErrEmptyPassword
	}

	salt := make([]byte, h.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("Argon2Hasher.Hash generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

	// Encode to standard PHC format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.memory, h.time, h.threads, b64Salt, b64Hash,
	)

	return encoded, nil
}

// Verify receives the user input password string and the encodedHash retrieved from the database
func (h *argon2Hasher) Verify(password, encodedHash string) (bool, error) {
	if strings.TrimSpace(password) == "" {
		return false, domain.ErrEmptyPassword
	}

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, domain.ErrInvalidPasswordFormat
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, domain.ErrInvalidPasswordFormat
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, domain.ErrInvalidPasswordFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("decode hash: %w", err)
	}

	// Recompute hash using stored parameters
	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(decodedHash)))

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1 {
		return true, nil
	}

	return false, nil
}
