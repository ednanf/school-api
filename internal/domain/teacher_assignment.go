package domain

import (
	"context"
	"time"
)

// TeacherAssignmentRepository defines the contract for database operations (in teacher_assignment_repo.go)
type TeacherAssignmentRepository interface {
	Create(ctx context.Context, t *TeacherAssignment) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*PopulatedTeacherAssignment, error)
	List(ctx context.Context, limit, offset int) ([]PopulatedTeacherAssignment, int, error)
	Update(ctx context.Context, id int, input PatchTeacherAssignmentInput) (*PopulatedTeacherAssignment, error)
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
	TeacherID *int `json:"teacher_id" db:"teacher_id" validate:"omitempty,required,gt=0"`
	ClassID   *int `json:"class_id" db:"class_id" validate:"omitempty,required,gt=0"`
	SubjectID *int `json:"subject_id" db:"subject_id" validate:"omitempty,required,gt=0"`
}

type TeacherSummary struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

// ClassSummary represents embedded class data inside an assignment
type ClassSummary struct {
	ID     int    `db:"id" json:"id"`
	Grade  int    `db:"grade" json:"grade"`
	Letter string `db:"letter" json:"letter"`
}

// SubjectSummary represents embedded subject data inside an assignment
type SubjectSummary struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// PopulatedTeacherAssignment is used for GET/List API responses
type PopulatedTeacherAssignment struct {
	ID        int            `db:"id" json:"id"`
	Teacher   TeacherSummary `db:"teacher" json:"teacher"`
	Class     ClassSummary   `db:"class" json:"class"`
	Subject   SubjectSummary `db:"subject" json:"subject"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}
