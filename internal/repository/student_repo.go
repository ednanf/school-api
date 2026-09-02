package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type studentRepo struct {
	db *sqlx.DB
}

// Constructor
func NewStudentRepository(db *sqlx.DB) domain.StudentRepository {
	return &studentRepo{db: db}
}

// TODO: Add count to `List` and `Batch` methods

func (r *studentRepo) Create(ctx context.Context, s *domain.Student) error {
	query := `
		INSERT INTO students (first_name, last_name, email, class, created_at, updated_at)
		VALUES (:first_name, :last_name, :email, :class, :created_at, :updated_at)
	`
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt

	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)
	return nil
}

func (r *studentRepo) GetByID(ctx context.Context, id int) (*domain.Student, error) {
	var s domain.Student
	query := "SELECT * FROM students WHERE id = ?"
	err := r.db.GetContext(ctx, &s, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) List(ctx context.Context, limit int, offset int) ([]domain.Student, error) {
	students := make([]domain.Student, 0)
	query := "SELECT * FROM students ORDER BY id LIMIT ? OFFSET ?"
	err := r.db.SelectContext(ctx, &students, query, limit, offset)
	return students, err
}

// TODO: Make the remaining methods
func (r *studentRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM students WHERE id = ?"

	// Execute query
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If 0 rows were affected, the ID did not exist in the db
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *studentRepo) Patch(ctx context.Context, s *domain.Student) error {
	return nil // Stub for now
}

func (r *studentRepo) BatchCreate(ctx context.Context, students []*domain.Student) error {
	return nil // Stub for now
}

func (r *studentRepo) BatchDelete(ctx context.Context, students []domain.Student) error {
	return nil // Stub for now
}
