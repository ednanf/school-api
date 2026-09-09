package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	// Initiate a caser to normalize capitalization
	caser := cases.Title(language.English)

	// Add timestamp and normalize values
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	s.Name = caser.String(s.Name)

	fmt.Printf("[DEBUG] s: %v\n", s)

	// Execute the db operation with `NamedExecContext` to match the named placeholders
	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return err
	}

	// Retrieve newly inserted entry's id to be able to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return err
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
		return err
	}

	// Check if any row was actually affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *subjectRepo) GetById(ctx context.Context, id int) (*domain.Subject, error) {
	var s domain.Subject

	query := "SELECT * FROM subjects WHERE id = ?"

	// Execute the db operation and assign it to the variable `s` if successful
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		// If the id is not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// Other errors
		return nil, err
	}

	return &s, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *subjectRepo) List(ctx context.Context, limit int, offset int) ([]domain.Subject, error) {
	// Make an empty slice to hold subjects
	subjects := make([]domain.Subject, 0)

	// Get all columns from table subjects, ordered by their ID, limited to a certain number
	query := "SELECT * FROM subjects ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	err := r.db.SelectContext(ctx, &subjects, query, limit, offset)

	// Return the results to be used
	return subjects, err
}

func (r *subjectRepo) Update(ctx context.Context, id int, input domain.PatchSubjectInput) (*domain.Subject, error) {
	// Fetch the current record from the db
	subject, err := r.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Overwrite only fields provided in the PATCH payload
	if input.Name != nil {
		subject.Name = *input.Name
	}
	subject.UpdatedAt = time.Now()

	query := `
		UPDATE subjects SET
			name = :name,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err = r.db.NamedExecContext(ctx, query, subject)
	if err != nil {
		return nil, err
	}

	return subject, nil
}
