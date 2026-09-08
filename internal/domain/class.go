package domain

import (
	"context"
	"time"
)

// ClassRepository defines the contract for database operations (in class_repo.go)
type ClassRepository interface {
	Create(ctx context.Context, class *Class) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Class, error)
	List(ctx context.Context, limit int, offset int) ([]Class, int, error)
	Update(ctx context.Context, id int, input PatchClassInput) (*Class, error)
}

// Class defines the shape of Class struct in the database
type Class struct {
	ID        int       `json:"id" db:"id"`
	Grade     int       `json:"grade" db:"grade" validate:"required,min=1,max=9"`
	Letter    string    `json:"letter" db:"letter" validate:"required,oneof=A B C D"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PatchClassInput defines the JSON payload for patching one class
type PatchClassInput struct {
	// Since the types are primitives, pointers must be used to avoid overwriting nil values. Change "required" to "omitempty" because the values are optional
	Grade  *int    `json:"grade" db:"grade" validate:"omitempty,min=1,max=9"`
	Letter *string `json:"letter" db:"letter" validate:"omitempty,oneof=A B C D"`
}
