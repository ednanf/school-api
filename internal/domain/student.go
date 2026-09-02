package domain

import (
	"context"
	"time"
)

// StudentRepository defines the contract for database operations (that happen in student_repo.go)
type StudentRepository interface {
	BatchCreate(ctx context.Context, students []Student) ([]Student, error)
	BatchDelete(ctx context.Context, ids []int) (int64, error)
	Create(ctx context.Context, student *Student) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Student, error)
	List(ctx context.Context, limit int, offset int) ([]Student, error)
	Update(ctx context.Context, id int, input PatchStudentInput) (*Student, error)
}

// Student defines the shape of Student struct in the database
type Student struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Class     string    `json:"class" db:"class" validate:"required,alphanum,max=20"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PatchStudentInput defines the JSON payload for inserting one student
type PatchStudentInput struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Class     *string `json:"class" validate:"omitempty"`
}

// BatchCreateInput defines the JSON payload for inserting multiple students
type BatchCreateInput struct {
	// `dive` tells validator to iterate into the slice and run field validation on each individual element
	Students []Student `json:"students" validate:"required,min=1,max=100,dive"`
}

// BatchDeleteInput defines the JSON payload for deleting multiple students by ID
type BatchDeleteInput struct {
	IDs []int `json:"ids" validate:"required,min=1,max=100,dive,gt=0"`
}
