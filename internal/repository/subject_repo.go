package repository

import (
	"context"
	"database/sql"
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

func (r *subjectRepo) GetById(ctx context.Context, id int) (*domain.Class, error) {
	return nil, nil
}

func (r *subjectRepo) List(ctx context.Context, limit int, offset int) ([]domain.Subject, int, error) {
	return nil, 0, nil
}

func (r *subjectRepo) Update(ctx context.Context, id int, input domain.PatchSubjectInput) (*domain.Subject, error) {
	return nil, nil
}
