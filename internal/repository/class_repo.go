package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

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

	// Add timestamp and normalize the casing
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Letter = strings.ToUpper(c.Letter)

	// Execute the db operation with `NamedExecContext` to match the named placeholders
	result, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return err
	}

	// Retrieve newly inserted entry's id to be able to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return err
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
		return err
	}

	// Check if any row was actually deleted
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If 0 rows were affected, the ID did not exist in the db
	if rowsaffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *classRepo) GetById(ctx context.Context, id int) (*domain.Class, error) {
	var c domain.Class

	query := "SELECT * FROM classes WHERE id = ?"

	// Execute the db operation and assign it to the variable `c` if successful
	if err := r.db.GetContext(ctx, &c, query, id); err != nil {
		// If the id is not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// Other errors
		return nil, err
	}

	return &c, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *classRepo) List(ctx context.Context, limit int, offset int) ([]domain.Class, error) {
	// Make an empty slice to hold classes
	classes := make([]domain.Class, 0)

	// Get all columns from the table classes, ordered by their ID, limited to a certain number
	query := "SELECT * FROM classes ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	err := r.db.SelectContext(ctx, &classes, query, limit, offset)

	// Return the results to be used
	return classes, err
}

func (r *classRepo) Update(ctx context.Context, id int, input domain.PatchClassInput) (*domain.Class, error) {
	// Fetch the current record from the db
	class, err := r.GetById(ctx, id)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if 404
	}

	// Overwrite only fields provided in the PATCH payload
	if input.Grade != nil {
		class.Grade = *input.Grade
	}
	if input.Letter != nil {
		class.Letter = *input.Letter
	}
	class.UpdatedAt = time.Now()

	query := `
		UPDATE classes SET
			grade = :grade,
			letter = :letter,
			updated_at = :updated_at
		WHERE id = :id
	`

	// Execute static SQL query using named placeholders
	_, err = r.db.NamedExecContext(ctx, query, class)
	if err != nil {
		return nil, err
	}

	return class, nil
}
