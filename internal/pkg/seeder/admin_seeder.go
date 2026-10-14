package seeder

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ednanf/school-api/internal/domain"
)

func SeedInitialAdmin(ctx context.Context, userRepo domain.UserRepository, hasher domain.PasswordHasher) error {
	// Check if a SUPER_ADMIN already exists
	exists, err := userRepo.ExistsByRole(ctx, "super_admin")
	if err != nil {
		return err
	}

	if exists {
		log.Println("[SYSTEM] Initial admin account already exists. Skipping seed.")
		return nil
	}

	// Read admin credentials from environment variables
	username := os.Getenv("INITIAL_ADMIN_USERNAME")
	email := os.Getenv("INITIAL_ADMIN_EMAIL")
	plainPassword := os.Getenv("INITIAL_ADMIN_PASSWORD")

	if username == "" || email == "" || plainPassword == "" {
		log.Println("[SYSTEM] Initial admin credentials not fully provided in env. Skipping seed.")
		return nil
	}

	// Hash password using Argon2
	hashedPassword, err := hasher.Hash(plainPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "super_admin",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save to DB
	if err := userRepo.Create(ctx, admin); err != nil {
		return err
	}

	log.Printf("[SYSTEM] Initial super_admin account created successfully (%s)", email)
	return nil
}
