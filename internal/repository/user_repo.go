package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (
			staff_id, username, email, password_hash,
			password_reset_token, password_token_expires,
			role, is_active, last_login_at, created_at, updated_at
		) VALUES (
			:staff_id, :username, :email, :password_hash,
			:password_reset_token, :password_token_expires,
			:role, :is_active, :last_login_at, :created_at, :updated_at
		)
	`

	result, err := r.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return fmt.Errorf("userRepo.Create execute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("userRepo.Create last insert id: %w", err)
	}

	u.ID = int(id)
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT
			id, staff_id, username, email, password_hash,
			password_reset_token, password_token_expires,
			role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE id = ?
	`
	var u domain.User
	if err := r.db.GetContext(ctx, &u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("userRepo.GetByID: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("userRepo.GetByID execute: %w", err)
	}
	return &u, nil
}

// GetByEmail is used interally for email lookup for authentication endpoints such as login
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT
			id, staff_id, username, email, password_hash,
			password_reset_token, password_token_expires,
			role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE email = ?
	`
	var u domain.User
	if err := r.db.GetContext(ctx, &u, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("userRepo.GetByEmail execute: %w", err)
	}
	return &u, nil
}

func (r *userRepo) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM users"); err != nil {
		return nil, 0, fmt.Errorf("userRepo.List count: %w", err)
	}

	if total == 0 {
		return []domain.User{}, 0, nil
	}

	query := `
		SELECT
			id, staff_id, username, email, password_hash,
			password_reset_token, password_token_expires,
			role, is_active, last_login_at, created_at, updated_at
		FROM users ORDER BY id ASC LIMIT ? OFFSET ?
	`
	users := make([]domain.User, 0, limit)
	if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("userRepo.List select: %w", err)
	}

	return users, total, nil
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET
			staff_id = :staff_id,
			username = :username,
			email = :email,
			password_hash = :password_hash,
			password_reset_token = :password_reset_token,
			password_token_expires = :password_token_expires,
			role = :role,
			is_active = :is_active,
			last_login_at = :last_login_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	result, err := r.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return fmt.Errorf("userRepo.Update execute: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepo.Update rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("userRepo.Update: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("userRepo.Delete execute: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepo.Delete rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("userRepo.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

// GetByPasswordResetToken finds a row that matches the temporary token. This is used in Auth service layer
func (r *userRepo) GetByPasswordResetToken(ctx context.Context, token string) (*domain.User, error) {
	query := `
		SELECT
			id, staff_id, username, email, password_hash,
			password_reset_token, password_token_expires,
			role, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE password_reset_token = ?
	`
	var u domain.User
	if err := r.db.GetContext(ctx, &u, query, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("userRepo.GetByPasswordResetToken: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("userRepo.GetByPasswordResetToken execute: %w", err)
	}
	return &u, nil
}
