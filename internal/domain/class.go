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
	List(ctx context.Context, limit int, offset int) ([]Class, error)
	Update(ctx context.Context, id int, input PatchClassInput) (*Class, error)
}

// Class defines the shape of Class struct in the database
type Class struct {
	ID        int       `json:"id" db:"id"`
	Grade     int       `json:"grade" db:"grade" validate:"required,min=1,max=8"`
	Letter    string    `json:"letter" db:"letter" validate:"required,oneof=A B"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PatchClassInput defines the JSON payload for inserting one class
type PatchClassInput struct {
	Grade  int    `json:"grade" db:"grade" validate:"required,min=1,max=8"`
	Letter string `json:"letter" db:"letter" validate:"required,oneof=A B"`
}
