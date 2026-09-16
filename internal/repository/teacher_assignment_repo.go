package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// taRepo stores the db connection and the repository methods attached to it
type taRepo struct {
	db *sqlx.DB
}

// NewTeacherAssignmentRepository receives a pointer to the database conneciton pool and returns a domain.TeacherAssignmentRepository
func NewTeacherAssignmentRepository(db *sqlx.DB) domain.TeacherAssignmentRepository {
	return &taRepo{db: db}
}

func (r *taRepo) Create(ctx context.Context, t *domain.TeacherAssignment) error {
	query := `
		INSERT INTO teacher_assignments (teacher_id, class_id, subject_id, created_at, updated_at)
		VALUES (:teacher_id, :class_id, :subject_id, :created_at, :updated_at)
	`

	// Add timestamp
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	result, err := r.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("taRepo.Create execute: %w", err)
	}

	// Grab newly inserted entry's id to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("taRepo.Create last inserted id: %w", err)
	}

	// Assign the id to the entry to show in the response
	t.ID = int(id)

	return nil
}

func (r *taRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM teacher_assignments WHERE id = ?"

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("taRepo.Delete execute: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taRepo.Delete rows affected: %w", err)
	}

	// If 0 rows were affected, the ID did not exist in the db
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *taRepo) GetById(ctx context.Context, id int) (*domain.TeacherAssignment, error) {
	return nil, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *taRepo) List(ctx context.Context, limit, offset int) ([]domain.TeacherAssignment, int, error) {
	return nil, 0, nil
}

func (r *taRepo) Update(ctx context.Context, id int, input domain.PatchTeacherAssignmentInput) (*domain.TeacherAssignment, error) {
	return nil, nil
}
