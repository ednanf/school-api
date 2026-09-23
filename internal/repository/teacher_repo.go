package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// teacherRepo stores the db connection and the repository methods attached to it
type teacherRepo struct {
	db *sqlx.DB
}

// NewTeacherRepository receives a pointer to the database connection pool and returns a domain.TeacherRepository, guaranteeing teacherRepo implements all required methods
func NewTeacherRepository(db *sqlx.DB) domain.TeacherRepository {
	return &teacherRepo{db: db}
}

// TODO: refactor to service layer

func (r *teacherRepo) Create(ctx context.Context, t *domain.Teacher) error {
	query := `
		INSERT INTO teachers (first_name, last_name, email, created_at, updated_at)
		VALUES (:first_name, :last_name, :email, :created_at, :updated_at)
	`

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
	query := "DELETE FROM teachers WHERE id = ?"

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("teacherRepo.Delete execute: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("teacherRepo.Delete rows affected: %w", err)
	}

	// If no rows were affected, the ID did not exist in the db
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *teacherRepo) GetById(ctx context.Context, id int) (*domain.Teacher, error) {
	var t domain.Teacher

	query := "SELECT id, first_name, last_name, email, is_active, created_at, updated_at FROM teachers WHERE id = ?"

	// Execute the query and assign it to the variable `t` if successful
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		// If the ID is not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		// Other errors
		return nil, fmt.Errorf("teacherRepo.GetByID execute: %w", err)
	}

	return &t, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *teacherRepo) List(ctx context.Context, limit int, offset int) ([]domain.Teacher, int, error) {
	// Make an empty slice to hold teachers
	teachers := make([]domain.Teacher, 0)

	// Get the total count across the entire table
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM teachers"
	if err := r.db.GetContext(ctx, &totalItems, countQuery); err != nil {
		return nil, 0, fmt.Errorf("teacherRepo.List count: %w", err)
	}

	// Get all columns from the table students, ordered by their ID, and limited to a certain number
	query := "SELECT id, first_name, last_name, email, is_active, created_at, updated_at FROM teachers ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	if err := r.db.SelectContext(ctx, &teachers, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("teacherRepo.List fetch: %w", err)
	}

	// Return results to be used
	return teachers, totalItems, nil
}

func (r *teacherRepo) Update(ctx context.Context, t *domain.Teacher) error {
	query := `
		UPDATE teachers SET
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("teacherRepo.Update execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("teacherService.Update rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
