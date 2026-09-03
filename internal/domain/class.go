package domain

import "time"

type Class struct {
	ID        int       `json:"id" db:"id"`
	Grade     int       `json:"grade" db:"grade" validate:"required,min=1,max=8"`
	Letter    string    `json:"letter" db:"letter" validate:"required,oneof=A B"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
