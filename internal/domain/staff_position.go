package domain

import (
	"context"
	"time"
)

type StaffPositionRepository interface {
	Create(ctx context.Context, s *StaffPosition) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*PopulatedStaffPosition, error)
	List(ctx context.Context, limit, offset int) (positions []PopulatedStaffPosition, totalItems int, err error)
	Update(ctx context.Context, id int) error
}

type StaffPositionService interface {
	Create(ctx context.Context, staffPosition *StaffPosition) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*PopulatedStaffPosition, error)
	List(ctx context.Context, limit, offset int) (positions []PopulatedStaffPosition, totalItems int, err error)
	Update(ctx context.Context, id int) (*StaffPosition, error)
}

// StaffPosition defines the shape of StaffPosition struct in the database
type StaffPosition struct {
	ID           int       `json:"id" db:"id"`
	DepartmentID int       `json:"department_id" db:"department_id" validate:"required,gt=0"`
	Title        string    `json:"title" db:"title" validate:"required,min=2,max=100"`
	Description  string    `json:"description" db:"description" validate:"required,min=2,max=300"`
	IsActive     bool      `json:"is_active" db:"is_active" validate:"required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type PatchStaffPosition struct {
	DepartmentID int    `json:"department_id" db:"department_id" validate:"omitempty,gt=0"`
	Title        string `json:"title" db:"title" validate:"omitempty,min=2,max=100"`
	Description  string `json:"description" db:"description" validate:"omitempty,min=2,max=300"`
	IsActive     bool   `json:"is_active" db:"is_active" validate:"omitempty"`
}

// DepartmentSummary represents embedded department data inside a staff position
type DepartmentSummary struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// PopulatedStaffPosition is used for GET/List API responses
type PopulatedStaffPosition struct {
	Department  DepartmentSummary `json:"department_id"`
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
}
