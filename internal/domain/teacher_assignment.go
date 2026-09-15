package domain

import (
	"context"
	"time"
)

// TeacherAssignmentRepository defines the contract for database operations (in teacher_assignment_repo.go)
type TeacherAssignmentRepository interface {
	Create(ctx context.Context, t *TeacherAssignment) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*TeacherAssignment, error)
	List(ctx context.Context, limit, offset int) ([]TeacherAssignment, int, error)
	Update(ctx context.Context, id int, input PatchTeacherAssignmentInput) (*TeacherAssignment, error)
}

// TeacherAssignment defines the shape of TeacherAssignment struct in the database
type TeacherAssignment struct {
	ID        int       `json:"id" db:"id"`
	TeacherID int       `json:"teacher_id" db:"teacher_id" validate:"required,gt=0"`
	ClassID   int       `json:"class_id" db:"class_id" validate:"required,gt=0"`
	SubjectID int       `json:"subject_id" db:"subject_id" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PatchTeacherAssignmentInput struct {
	TeacherID int `json:"teacher_id" db:"teacher_id" validate:"required,gt=0"`
	ClassID   int `json:"class_id" db:"class_id" validate:"required,gt=0"`
	SubjectID int `json:"subject_id" db:"subject_id" validate:"required,gt=0"`
}
