package domain

import (
	"context"
	"strings"
	"time"

	"golang.org/x/text/cases"
)

// TeacherRepository defines the contract for the database operations (in teacher_repo.go)
type TeacherRepository interface {
	Create(ctx context.Context, t *Teacher) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Teacher, error)
	List(ctx context.Context, limit int, offset int) ([]Teacher, int, error)
	Update(ctx context.Context, id int, input PatchTeacherInput) (*Teacher, error)
}

// Teacher defines the shape of Student struct in the database
type Teacher struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PatchTeacherInput defines the JSON payload for inserting one student
type PatchTeacherInput struct {
	// Since the types are primitives, pointers must be used to avoid overwriting nil values. Change "required" to "omitempty" because the values are optional
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Email     *string `json:"email" validate:"omitempty,email"`
	ClassID   *int    `json:"class_id" validate:"omitempty,gt=0"`
}

// Normalize applies string normalization to user input fields
func (t *Teacher) Normalize(caser cases.Caser) {
	t.FirstName = caser.String(strings.ToLower(strings.TrimSpace(t.FirstName)))
	t.LastName = caser.String(strings.ToLower(strings.TrimSpace(t.LastName)))
	t.Email = strings.ToLower(strings.TrimSpace(t.Email))
}
