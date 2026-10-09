package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims defines the payload baked into the JWT
type CustomClaims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}

// TokenService abstracts token generation and verification
type TokenService interface {
	GenerateToken(user *User, ttl time.Duration) (string, error)
	ValidateToken(tokenStr string) (*CustomClaims, error)
}

// AuthService contracts the authentication business logic
type AuthService interface {
	Login(ctx context.Context, input LoginInput) (*User, string, error)
	ChangePassword(ctx context.Context, userID int, input ChangePasswordInput) error
	ForgotPassword(ctx context.Context, input ForgotPasswordInput) (string, error)
	ResetPassword(ctx context.Context, input ResetPasswordInput) error
}

// DTOs for Auth Endpoints
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse uses the existing UserResponse from domain/user.go
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
