package domain

import (
	"context"
	"strings"
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleStaff   UserRole = "staff"
	RoleTeacher UserRole = "teacher"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error) // Used internally in the login service
	List(ctx context.Context, limit, offset int) ([]User, int, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int) error
}

type UserService interface {
	Create(ctx context.Context, u *User, plainPassword string) error
	GetByID(ctx context.Context, id int) (*User, error)
	List(ctx context.Context, page, limit int) ([]User, int, int, int, error)
	Update(ctx context.Context, id int, input PatchUserInput) (*User, error)
	Delete(ctx context.Context, id int) error
}

// User represents both the HTTP payload/response and the Database model.
// Plaintext password is taken as a separate parameter in Service/HTTP layers
// so it never needs a temporary struct field.
type User struct {
	ID                   int        `json:"id" db:"id"`
	StaffID              *int       `json:"staff_id" db:"staff_id"`
	Username             string     `json:"username" db:"username"`
	Email                string     `json:"email" db:"email"`
	PasswordHash         string     `json:"-" db:"password_hash"`
	PasswordResetToken   *string    `json:"-" db:"password_reset_token"`
	PasswordTokenExpires *time.Time `json:"-" db:"password_token_expires"`
	Role                 UserRole   `json:"role" db:"role"`
	IsActive             bool       `json:"is_active" db:"is_active"`
	LastLoginAt          *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateUserInput carries payload fields for POST /users (including plain password for validation)
type CreateUserInput struct {
	StaffID  *int     `json:"staff_id" validate:"omitempty"`
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required,oneof=admin staff teacher"`
	IsActive bool     `json:"is_active" validate:"boolean"`
}

// PatchUserInput handles optional updates for PATCH /users/{id}
type PatchUserInput struct {
	StaffID  *int      `json:"staff_id" validate:"omitempty"`
	Username *string   `json:"username" validate:"omitempty,min=3,max=50"`
	Email    *string   `json:"email" validate:"omitempty,email"`
	Role     *UserRole `json:"role" validate:"omitempty,oneof=admin staff teacher"`
	IsActive *bool     `json:"is_active" validate:"omitempty"`
}

// UserResponse carries payload fields for the response of POST/PATCH /users, ensuring no password/hash is present in the struct
type UserResponse struct {
	ID       int      `json:"id"`
	StaffID  *int     `json:"staff_id,omitempty"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
	IsActive bool     `json:"is_active"`
}

func (p PatchUserInput) HasUpdates() bool {
	return p.StaffID != nil || p.Username != nil || p.Email != nil || p.Role != nil || p.IsActive != nil
}

func (u *User) Normalize() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
}
