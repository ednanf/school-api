package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
)

type userService struct {
	repo   domain.UserRepository
	hasher domain.PasswordHasher
}

func NewUserService(repo domain.UserRepository, hasher domain.PasswordHasher) domain.UserService {
	return &userService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *userService) Create(ctx context.Context, u *domain.User, plainPassword string) error {
	u.Normalize()

	// Use Argon2Hasher via domain.PasswordHasher interface
	hashedPassword, err := s.hasher.Hash(plainPassword)
	if err != nil {
		return fmt.Errorf("userService.Create hash password: %w", err)
	}
	u.PasswordHash = hashedPassword

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	if err := s.repo.Create(ctx, u); err != nil {
		return fmt.Errorf("userService.Create: %w", err)
	}

	return nil
}

func (s *userService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userService.GetByID: %w", err)
	}
	return user, nil
}

func (s *userService) List(ctx context.Context, page, limit int) (items []domain.User, totalItems, pageNum, limitNum int, err error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, totalItems, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("userService.List: %w", err)
	}

	return users, totalItems, page, limit, nil
}

func (s *userService) Update(ctx context.Context, id int, input domain.PatchUserInput) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userService.Update fetch: %w", err)
	}

	if input.StaffID != nil {
		user.StaffID = input.StaffID
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	user.Normalize()
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("userService.Update execute: %w", err)
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("userService.Delete: %w", err)
	}
	return nil
}
