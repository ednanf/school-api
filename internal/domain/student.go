package domain

import (
	"context"
	"strings"
	"time"

	"golang.org/x/text/cases"
)

// StudentRepository defines the contract for database operations (in student_repo.go)
type StudentRepository interface {
	BatchCreate(ctx context.Context, students []Student) ([]Student, int, error)
	BatchDelete(ctx context.Context, ids []int) (int64, error)
	BatchUpdate(ctx context.Context, updates []BatchUpdateStudentItem) ([]Student, int, error)
	BatchUpdateClass(ctx context.Context, ids []int, classID int) (int64, error)
	Create(ctx context.Context, student *Student) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Student, error)
	List(ctx context.Context, limit int, offset int) ([]Student, int, error)
	Update(ctx context.Context, id int, input PatchStudentInput) (*Student, error)
}

// Student defines the shape of Student struct in the database
type Student struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	ClassID   int       `json:"class_id" db:"class_id" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PatchStudentInput defines the JSON payload for inserting one student
type PatchStudentInput struct {
	// Since the types are primitives, pointers must be used to avoid overwriting nil values. Change "required" to "omitempty" because the values are optional
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Email     *string `json:"email" validate:"omitempty,email"`
	ClassID   *int    `json:"class_id" validate:"omitempty,gt=0"`
}

// BatchCreateStudentInput defines the JSON payload for inserting multiple students
type BatchCreateStudentInput struct {
	// `dive` tells validator to iterate into the slice and run field validation on each individual element
	Students []Student `json:"students" validate:"required,min=1,max=2000,dive"`
}

// BatchDeleteStudentInput defines the JSON payload for deleting multiple students by ID
type BatchDeleteStudentInput struct {
	IDs []int `json:"ids" validate:"required,min=1,max=100,dive,gt=0"`
}

// BatchUpdateStudentItem is a single item inside the batch update payload
type BatchUpdateStudentItem struct {
	ID int `json:"id" validate:"required,gt=0"`
	PatchStudentInput
}

// BatchUpdateStudentInput is a wrapper DTO for batch update endpoint
type BatchUpdateStudentInput struct {
	Students []BatchUpdateStudentItem `json:"students" validate:"required,min=1,max=1000,dive"`
}

// Normalize applies string normalization to user input fields
func (s *Student) Normalize(caser cases.Caser) {
	s.FirstName = caser.String(strings.ToLower(strings.TrimSpace(s.FirstName)))
	s.LastName = caser.String(strings.ToLower(strings.TrimSpace(s.LastName)))
	s.Email = strings.ToLower(strings.TrimSpace(s.Email))
}
