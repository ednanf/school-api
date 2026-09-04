package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// studentRepo stores the db connnection and the repository methods attached to it
type studentRepo struct {
	db *sqlx.DB
}

// NewStudentRepository receives a pointer to the database connection pool and returns a domain.StudentRepository, guaranteeing studentRepo implements all required methods
func NewStudentRepository(db *sqlx.DB) domain.StudentRepository {
	return &studentRepo{db: db}
}

// TODO: Add count to `List` and `Batch` methods

func (r *studentRepo) Create(ctx context.Context, s *domain.Student) error {
	query := `
        INSERT INTO students (first_name, last_name, email, class_id, created_at, updated_at)
        VALUES (:first_name, :last_name, :email, :class_id, :created_at, :updated_at)
    `
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

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
	// Make an empty slice to hold students
	students := make([]domain.Student, 0)

	// Get all columns from the table students, ordered by their ID, and limited to a certain number
	query := "SELECT * FROM students ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	err := r.db.SelectContext(ctx, &students, query, limit, offset)

	// Return the results to be used
	return students, err
}

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

func (r *studentRepo) Update(ctx context.Context, id int, input domain.PatchStudentInput) (*domain.Student, error) {
	// Fetch current record from DB
	student, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if 404
	}

	// Overwrite only fields provided in the PATCH payload
	if input.FirstName != nil {
		student.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		student.LastName = *input.LastName
	}
	if input.Email != nil {
		student.Email = *input.Email
	}
	if input.ClassID != nil {
		student.ClassID = *input.ClassID
	}
	student.UpdatedAt = time.Now()

	// Execute static SQL query using sqlx named placeholders
	query := `
        UPDATE students SET
            first_name = :first_name,
            last_name = :last_name,
            email = :email,
            class_id = :class_id,
            updated_at = :updated_at
        WHERE id = :id
    `

	_, err = r.db.NamedExecContext(ctx, query, student)
	if err != nil {
		return nil, err
	}

	return student, nil
}

// BatchCreate accepts 100 entries at most
func (r *studentRepo) BatchCreate(ctx context.Context, students []domain.Student) ([]domain.Student, error) {
	if len(students) == 0 {
		return students, nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	query := `
        INSERT INTO students (first_name, last_name, email, class_id, created_at, updated_at)
        VALUES (:first_name, :last_name, :email, :class_id, :created_at, :updated_at)
    `

	for i := range students {
		students[i].CreatedAt = now
		students[i].UpdatedAt = now

		res, err := tx.NamedExecContext(ctx, query, students[i])
		if err != nil {
			return nil, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		students[i].ID = int(id)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return students, nil
}

// BatchDelete accepts 100 entries at most
func (r *studentRepo) BatchDelete(ctx context.Context, ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Expand the slice into positional placeholders: WHERE id IN (?, ?, ...)
	query, args, err := sqlx.In("DELETE FROM students WHERE id IN (?)", ids)
	if err != nil {
		return 0, err
	}

	// Rebind query to match MariaDB driver syntax
	query = r.db.Rebind(query)

	// Execute query
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	// Return total number of deleted rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
