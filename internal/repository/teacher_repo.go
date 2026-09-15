package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// teacherRepo stores the db connection and the repository methods attached to it
type teacherRepo struct {
	db *sqlx.DB
}

// NewTeacherRepository receives a pointer to the database connection pool and returns a domain.TeacherRepository, guaranteeing teacherRepo implements all required methods
func NewTeacherRepository(db *sqlx.DB) domain.TeacherRepository {
	return &teacherRepo{db: db}
}

func (r *teacherRepo) Create(ctx context.Context, t *domain.Teacher) error {
	query := `
		INSERT INTO teachers (first_name, last_name, email, created_at, updated_at)
		VALUES (:first_name, :last_name, :email, :created_at, :updated_at)
	`

	// Initialize caser to normalize capitalization
	caser := cases.Title(language.English)

	// Normalize fields explicitly (Trim + Case handling)
	t.Normalize(caser)

	// Add timestamp
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	// Execute the db operation with `NamedExecContext` to match the named placeholders
	result, err := r.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("teacherRepo.Create execute: %w", err)
	}

	// Grab the newly inserted entry's id to be able to send a complete response
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("studentRepo.Create last inserted id: %w", err)
	}

	// Assign the received id to the entry in order to show in the response
	t.ID = int(id)

	return nil
}

func (r *teacherRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *teacherRepo) GetById(ctx context.Context, id int) (*domain.Teacher, error) {
	return nil, nil
}

func (r *teacherRepo) List(ctx context.Context, limit int, offset int) ([]domain.Teacher, int, error) {
	return nil, 0, nil
}

func (r *teacherRepo) Update(ctx context.Context, id int, input domain.PatchTeacherInput) (*domain.Teacher, error) {
	return nil, nil
}
