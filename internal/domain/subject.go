package domain

import (
	"context"
	"time"
)

// SubjectRepository defines the contract for database operations (in subject_repo.go)
type SubjectRepository interface {
	Create(ctx context.Context, s *Subject) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Class, error)
	List(ctx context.Context, limit int, offset int) ([]Subject, int, error)
	Update(ctx context.Context, id int, input PatchSubjectInput) (*Subject, error)
}

// Subject defines the shape of Subject struct in the database
type Subject struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=20"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PatchSubjectInput defines the JSON payload for patching one subject
type PatchSubjectInput struct {
	Name *string `json:"name" db:"name" validate:"required,min=1,max=20"` // Required since it's the only patchable field
}
