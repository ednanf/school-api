package domain

import (
	"context"
	"strings"
	"time"

	"golang.org/x/text/cases"
)

type DepartmentRepository interface {
	Create(ctx context.Context, d Department) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Department, error)
	List(ctx context.Context, limit int, offset int) ([]Department, int, error)
	Update(ctx context.Context, id int, input PatchDepartmentInput) (*Department, error)
}

// Department defines the shape of Department struct in the database
type Department struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=50"`
	Code        string    `json:"code" db:"code" validate:"omitempty,min=2,max=50"`
	Description string    `json:"description" db:"description" validate:"required,min=2,max=200"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PatchDepartmentInput defines the JSON payload for inserting a department
type PatchDepartmentInput struct {
	Name        string `json:"name" db:"name" validate:"required,min=2,max=50"`
	Code        string `json:"code" db:"code" validate:"omitempty,min=2,max=50"`
	Description string `json:"description" db:"description" validate:"required,min=2,max=200"`
}

// TODO: Verify if code is outputting uppercase only

// Normalize applies string normalization to user input fields
func (s *Department) Normalize(caser cases.Caser) {
	s.Name = caser.String(strings.ToLower(strings.TrimSpace(s.Name)))
	s.Code = caser.String(strings.ToUpper(strings.TrimSpace(s.Code)))
	s.Description = caser.String(strings.ToLower(strings.TrimSpace(s.Description)))
}
