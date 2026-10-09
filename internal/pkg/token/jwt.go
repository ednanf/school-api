package token

import (
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secretKey string, issuer string) domain.TokenService {
	return &jwtService{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

func (j *jwtService) GenerateToken(user *domain.User, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := domain.CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,

		// Registered claims
		Issuer:    j.issuer,
		Subject:   fmt.Sprintf("%d", user.ID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("JWTService.GenerateToken sign: %w", err)
	}

	return signedToken, nil
}

func (j *jwtService) ValidateToken(tokenStr string) (*domain.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&domain.CustomClaims{},
		func(t *jwt.Token) (any, error) {
			// Validate signing method to prevent algorithm downgrade attacks
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return j.secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("JWTService.ValidateToken parse: %w: %w", domain.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*domain.CustomClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}
