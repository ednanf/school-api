package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// subjectRepo stores the db connection and has the repository methods attached to it
type subjectRepo struct {
	db *sqlx.DB
}

// NewSubjectRepository receives a pointer to the database connection pool and returns a domain.SubjectRepository, guaranteeing subjectRepo implements all required methods
func NewSubjectRepository(db *sqlx.DB) domain.SubjectRepository {
	return &subjectRepo{db: db}
}

func (r *subjectRepo) Create(ctx context.Context, s *domain.Subject) error {
	query := `
		INSERT INTO subjects (name, created_at, updated_at)
		VALUES(:name, :created_at, :updated_at)
	`

	// Execute the db operation with `NamedExecContext` to match the named placeholders
	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("subjectRepo.Create execute: %w", err)
	}

	// Retrieve newly inserted entry's id to be able to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("subjectRepo.Create last insert id: %w", err)
	}

	// Assign the received id to the entry in order to show in the response
	s.ID = int(id)

	return nil
}

func (r *subjectRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM subjects WHERE id = ?"

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("subjectRepo.Delete execute: %w", err)
	}

	// Check if any row was actually affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("subjectRepo.Delete rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *subjectRepo) GetById(ctx context.Context, id int) (*domain.Subject, error) {
	var s domain.Subject

	query := "SELECT id, name, is_active, created_at, updated_at FROM subjects WHERE id = ?"

	// Execute the db operation and assign it to the variable `s` if successful
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		// If the id is not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		// Other errors
		return nil, fmt.Errorf("subjectRepo.GetById execute: %w", err)
	}

	return &s, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *subjectRepo) List(ctx context.Context, limit int, offset int) ([]domain.Subject, int, error) {
	// Make an empty slice to hold subjects
	subjects := make([]domain.Subject, 0)

	// Get the total count across the entire table
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM subjects"
	if err := r.db.GetContext(ctx, &totalItems, countQuery); err != nil {
		return nil, 0, fmt.Errorf("subjectRepo.List count: %w", err)
	}

	// Get all columns from table subjects, ordered by their ID, limited to a certain number
	query := "SELECT id, name, is_active, created_at, updated_at FROM subjects ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	if err := r.db.SelectContext(ctx, &subjects, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("subjectRepo.List fetch: %w", err)
	}

	// Return the results to be used
	return subjects, totalItems, nil
}

func (r *subjectRepo) Update(ctx context.Context, s *domain.Subject) error {
	query := `
		UPDATE subjects SET
			name = :name,
			updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("subjectRepo.Update exec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("subjectRepo.Update rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
