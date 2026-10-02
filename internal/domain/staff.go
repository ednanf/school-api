package domain

import (
	"context"
	"strings"
	"time"

	"golang.org/x/text/cases"
)

type StaffRepository interface {
	Create(ctx context.Context, s *StaffPosition) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*PopulatedStaff, error)
	List(ctx context.Context, page, offset int) (employees []PopulatedStaff, totalItems int, err error)
	Update(ctx context.Context, s *Staff) error
}

type StaffService interface {
	Create(ctx context.Context, e *Staff) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*PopulatedStaff, error)
	List(ctx context.Context, limit, offset int) (employees []PopulatedStaff, totalItems int, err error)
	Update(ctx context.Context, id int, input PatchStaffInput) (*Staff, error)
}

type Staff struct {
	ID         int       `json:"id" db:"id"`
	PositionID int       `json:"position_id" db:"position_id" validate:"required,gt=0"`
	FirstName  string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName   string    `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Email      string    `json:"email" db:"email" validate:"required,email"`
	HireDate   time.Time `json:"hire_date" db:"hire_date" validate:"required"`
	IsActive   bool      `json:"is_active" db:"is_active" validate:"required"`
	Created_At time.Time `json:"created_at" db:"created_at"`
	Updated_At time.Time `json:"updated_at" db:"updated_at"`
}

// StaffPositionSummary represents embedded staff position data inside PopulatedStaff
type StaffPositionSummary struct {
	ID          int               `json:"id" db:"id"`
	Department  DepartmentSummary `json:"department" db:"department"`
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description" db:"description"`
}

// TODO: Verify if the Position type is properly setup for chaining a JOIN operation
type PopulatedStaff struct {
	ID         int                  `json:"id" db:"id"`
	Position   StaffPositionSummary `json:"position" db:"position"`
	FirstName  string               `json:"first_name" db:"first_name"`
	LastName   string               `json:"last_name" db:"last_name"`
	Email      string               `json:"email" db:"email"`
	HireDate   time.Time            `json:"hire_date"`
	IsActive   bool                 `json:"is_active"`
	Created_At time.Time            `json:"created_at" db:"created_at"`
	Updated_At time.Time            `json:"updated_at" db:"updated_at"`
}

type PatchStaffInput struct {
	PositionID *int       `json:"position_id" db:"position_id" validate:"omitempty,gt=0"`
	FirstName  *string    `json:"first_name" db:"first_name" validate:"omitempty,min=2,max=50"`
	LastName   *string    `json:"last_name" db:"last_name" validate:"omitempty,min=2,max=50"`
	Email      string     `json:"email" db:"email" validate:"omitempty,email"`
	HireDate   *time.Time `json:"hire_date" db:"hire_date" validate:"omitempty"`
	IsActive   *bool      `json:"is_active" db:"is_active" validate:"omitempty"`
}

// HasUpdates returns true if at least one field is provided in the patch payload
func (p PatchStaffInput) HasUpdates() bool {
	return p.PositionID != nil || p.FirstName != nil || p.LastName != nil || p.HireDate != nil || p.IsActive != nil
}

func (s *Staff) Normalize(caser cases.Caser) {
	s.FirstName = caser.String(strings.ToLower(strings.TrimSpace(s.FirstName)))
	s.LastName = caser.String(strings.ToLower(strings.TrimSpace(s.LastName)))
	s.Email = strings.ToLower(strings.TrimSpace(s.Email))
}
