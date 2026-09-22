package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// classRepo stores the db connnection and has the repository methods attached to it
type classRepo struct {
	db *sqlx.DB
}

// NewClassRepository receives a pointer to the database connection pool and returns a domain.ClassRepository, guaranteeing classRepo implements all required methods
func NewClassRepository(db *sqlx.DB) domain.ClassRepository {
	return &classRepo{db: db}
}

func (r *classRepo) Create(ctx context.Context, c *domain.Class) error {
	query := `
		INSERT INTO classes (grade, letter, created_at, updated_at)
		VALUES (:grade, :letter, :created_at, :updated_at)
	`

	// Execute the db operation with `NamedExecContext` to match the named placeholders
	result, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("classRepo.Create execute: %w", err)
	}

	// Retrieve newly inserted entry's id to be able to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("classRepo.Create last insert id: %w", err)
	}

	// Assign the received id to the entry in order to show in the response
	c.ID = int(id)

	return nil
}

func (r *classRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM classes WHERE id = ?"

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("classRepo.Delete execute: %w", err)
	}

	// Check if any row was actually deleted
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepo.Delete rows affected: %w", err)
	}

	// If 0 rows were affected, the ID did not exist in the db
	if rowsaffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *classRepo) GetById(ctx context.Context, id int) (*domain.Class, error) {
	var c domain.Class

	query := "SELECT id, grade, letter, is_active, created_at, updated_at FROM classes WHERE id = ?"

	// Execute the db operation and assign it to the variable `c` if successful
	if err := r.db.GetContext(ctx, &c, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("classRepo.GetById execute: %w", err)
	}

	return &c, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *classRepo) List(ctx context.Context, limit int, offset int) ([]domain.Class, int, error) {
	// Make an empty slice to hold classes
	classes := make([]domain.Class, 0)

	// Get the total count across the entire table
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM classes"
	if err := r.db.GetContext(ctx, &totalItems, countQuery); err != nil {
		return nil, 0, fmt.Errorf("classRepo.List count: %w", err)
	}

	// Get all columns from the table classes, ordered by their ID, limited to a certain number
	query := "SELECT id, grade, letter, is_active, created_at, updated_at FROM classes ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	err := r.db.SelectContext(ctx, &classes, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("classRepo.List fetch: %w", err)
	}

	// Return the results to be used
	return classes, totalItems, nil
}

// ListStudentsByClassId retrieves all students assigned to a specific class ID
func (r *classRepo) ListStudentsByClassId(ctx context.Context, classID int, limit int, offset int) ([]domain.Student, int, error) {
	students := make([]domain.Student, 0)

	// Count total students belonging strictly to this class ID
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM students WHERE class_id = ?"
	if err := r.db.GetContext(ctx, &totalItems, countQuery, classID); err != nil {
		return nil, 0, fmt.Errorf("classRepo.ListStudentsByClassId count: %w", err)
	}

	query := `
		SELECT id, first_name, last_name, email, class_id, is_active, created_at, updated_at
		FROM students
		WHERE class_id = ?
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	if err := r.db.SelectContext(ctx, &students, query, classID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("classRepo.ListStudentsByClassId fetch: %w", err)
	}

	return students, totalItems, nil
}

func (r *classRepo) Update(ctx context.Context, c *domain.Class) error {
	query := `
		UPDATE classes SET
			grade = :grade,
			letter = :letter,
			updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("classRepo.Update execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepo.Update rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
