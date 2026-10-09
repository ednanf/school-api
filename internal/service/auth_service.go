package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"strings"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	logger "github.com/ednanf/school-api/internal/pkg/loggers"
)

type authService struct {
	userRepo domain.UserRepository
	hasher   domain.PasswordHasher
	tokens   domain.TokenService
	mailer   domain.Mailer
	jwtTTL   time.Duration
}

func NewAuthService(
	userRepo domain.UserRepository,
	hasher domain.PasswordHasher,
	tokens domain.TokenService,
	mailer domain.Mailer,
	jwtTTL time.Duration,
) domain.AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
		tokens:   tokens,
		mailer:   mailer,
		jwtTTL:   jwtTTL,
	}
}

// Auth shares the functionality of the User domain, therefore, the Auth service layer works mostly with the User repository/service layers

func (s *authService) Login(ctx context.Context, input domain.LoginInput) (user *domain.User, token string, err error) {

	// Normalize email
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Fetch user by email via repository
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, "", domain.ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("authService.Login fetch user: %w", err)
	}

	// Check if active
	if !user.IsActive {
		return nil, "", domain.ErrUserInactive
	}

	// Verify password hash using Argon2
	isValid, err := s.hasher.Verify(input.Password, user.PasswordHash)
	if err != nil || !isValid { // TODO: check if the !isValid is not inverted here
		return nil, "", domain.ErrInvalidCredentials
	}

	// Update last_login_at timestamp
	now := time.Now().UTC()
	user.LastLoginAt = &now

	// Persist changes via the user repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, "", fmt.Errorf("authService.Login update last login: %w", err)
	}

	// Generate JWT token
	token, err = s.tokens.GenerateToken(user, s.jwtTTL)
	if err != nil {
		return nil, "", fmt.Errorf("authService.Login generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID int, input domain.ChangePasswordInput) error {

	// Fetch user by email via user repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("authService.ChangePassword fetch user: %w", err)
	}

	// Verify current password hash against the incoming input password
	valid, err := s.hasher.Verify(input.CurrentPassword, user.PasswordHash)
	if err != nil || !valid {
		return domain.ErrInvalidCredentials
	}

	// Hash new password
	newHash, err := s.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("authService.ChangePassword hash: %w", err)
	}

	user.PasswordHash = newHash
	user.PasswordResetToken = nil
	user.PasswordTokenExpires = nil
	user.UpdatedAt = time.Now().UTC()

	// Persist the change via the user repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("authService.ChangePassword save user: %w", err)
	}

	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, input domain.ForgotPasswordInput) (string, error) {

	// Normalize email
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// To prevent user enumeration attacks, return without error
			return "", nil
		}
		return "", fmt.Errorf("authService.ForgotPassword fetch user: %w", err)
	}

	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("authService.ForgotPassword generate token: %w", err)
	}
	resetToken := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().UTC().Add(1 * time.Hour)

	user.PasswordResetToken = &resetToken
	user.PasswordTokenExpires = &expiresAt
	user.UpdatedAt = time.Now().UTC()

	// Persist changes via user repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", fmt.Errorf("authService.ForgotPassword update token: %w", err)
	}

	if err := s.mailer.SendPasswordResetEmail(ctx, user.Email, resetToken); err != nil {
		// Log email failure but don't crash the request
		logger.DebugLogger(err)
	}

	// Return resetToken so your notification service or caller can deliver it.
	return resetToken, nil
}

func (s *authService) ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error {
	// Fetch user by reset token
	user, err := s.userRepo.GetByPasswordResetToken(ctx, input.Token)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrInvalidToken
		}
		return fmt.Errorf("authService.ResetPassword fetch user: %w", err)
	}

	// Validate token expiration
	if user.PasswordTokenExpires == nil || time.Now().UTC().After(*user.PasswordTokenExpires) {
		return domain.ErrInvalidToken
	}

	// Hash new password
	newHash, err := s.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("authService.ResetPassword hash: %w", err)
	}

	// Update user state & clear reset token
	user.PasswordHash = newHash
	user.PasswordResetToken = nil
	user.PasswordTokenExpires = nil
	user.UpdatedAt = time.Now().UTC()

	// Persist changes via user repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("authService.ResetPassword update user: %w", err)
	}

	return nil
}
