package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// BulkCreate receives a context and a slice of type Student. It accepts 1000 entries at most. It returns a slice containing all created students, the total amount of created entries and errors
func (r *studentRepo) BulkCreate(ctx context.Context, students []domain.Student) ([]domain.Student, int, error) {
	if len(students) == 0 {
		return students, 0, nil
	}

	// Initiate a db transaction and defer a rollback in case of errors
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("studentRepo.BulkCreate begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO students (first_name, last_name, email, class_id, is_active, created_at, updated_at)
		VALUES (:first_name, :last_name, :email, :class_id, :is_active, :created_at, :updated_at)
	`

	total := 0

	// Loop through the students slice and execute the queries
	for i := range students {
		res, err := tx.NamedExecContext(ctx, query, students[i])
		if err != nil {
			return nil, 0, fmt.Errorf("studentRepo.BulkCreate insert at index %d: %w", i, err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return nil, 0, fmt.Errorf("studentRepo.BulkCreate last insert id at index %d: %w", i, err)
		}

		students[i].ID = int(id)
		total++
	}

	// Commit the trasaction to the database
	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("studentRepo.BulkCreate commit: %w", err)
	}

	return students, total, nil
}

// BulkDelete receives a context and a slice of ids of type int. It accepts a maximum of 100 ids. It returns a total and errors
func (r *studentRepo) BulkDelete(ctx context.Context, ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Expand slice into dynamic IN placeholders: WHERE id IN (?, ?, ...)
	query, args, err := sqlx.In("DELETE FROM students WHERE id IN (?)", ids)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkDelete query build: %w", err)
	}

	// Rebind query to driver syntax (?)
	query = r.db.Rebind(query)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkDelete execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkDelete rows affected: %w", err)
	}

	return rowsAffected, nil
}

// BulkUpdate receives a slice of batch inputs and updates each student record in a single transaction.
func (r *studentRepo) BulkUpdate(ctx context.Context, students []domain.Student) (int64, error) {
	if len(students) == 0 {
		return 0, nil
	}

	// Start a transaction with eventual rollbacks in case of errors
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdate begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE students SET
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			class_id = :class_id,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`

	var totalRowsAffected int64

	// Loop over the slice, saving the value of rows affected outside for future use
	for i := range students {
		// Expect students[i] to already be merged, normalized, and stamped by the Service layer
		res, err := tx.NamedExecContext(ctx, query, students[i])
		if err != nil {
			return 0, fmt.Errorf("studentRepo.BulkUpdate execute at index %d (ID %d): %w", i, students[i].ID, err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("studentRepo.BulkUpdate rows affected at index %d: %w", i, err)
		}

		if rows == 0 {
			return 0, fmt.Errorf("studentRepo.BulkUpdate student ID %d not found: %w", students[i].ID, domain.ErrNotFound)
		}

		totalRowsAffected += rows
	}

	// Commit the changes to the database
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdate commit: %w", err)
	}

	return totalRowsAffected, nil
}

// BulkUpdateClass reassigns a slice of student IDs to a new class_id in a single execution.
func (r *studentRepo) BulkUpdateClass(ctx context.Context, ids []int, classID int, updatedAt time.Time) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	rawQuery := "UPDATE students SET class_id = ?, updated_at = ? WHERE id IN (?)"

	// Expand slice into dynamic IN placeholders: WHERE id IN (?, ?, ...)
	query, args, err := sqlx.In(rawQuery, classID, updatedAt, ids)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdateClass query build: %w", err)
	}

	// Rebind query to match driver syntax (?)
	query = r.db.Rebind(query)

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdateClass execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdateClass rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (r *studentRepo) Create(ctx context.Context, s *domain.Student) error {
	query := `
		INSERT INTO students (first_name, last_name, email, class_id, is_active, created_at, updated_at)
		VALUES (:first_name, :last_name, :email, :class_id, :is_active, :created_at, :updated_at)
	`

	// Expect s to already be normalized and stamped by the Service layer
	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("studentRepo.Create execute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("studentRepo.Create last insert id: %w", err)
	}

	s.ID = int(id)
	return nil
}

func (r *studentRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM students WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("studentRepo.Delete execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("studentRepo.Delete rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("studentRepo.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *studentRepo) GetByID(ctx context.Context, id int) (*domain.PopulatedStudent, error) {
	query := `
		SELECT
			s.id, s.first_name, s.last_name, s.email, s.is_active, s.created_at, s.updated_at,
			c.id AS "class.id",
			c.grade AS "class.grade",
			c.letter AS "class.letter"
		FROM students s
		INNER JOIN classes c ON s.class_id = c.id
		WHERE s.id = ?
	`

	var student domain.PopulatedStudent
	if err := r.db.GetContext(ctx, &student, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("studentRepo.GetByID: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("studentRepo.GetByID execute: %w", err)
	}

	return &student, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *studentRepo) List(ctx context.Context, limit, offset int) ([]domain.PopulatedStudent, int, error) {
	countQuery := "SELECT COUNT(*) FROM students"
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("studentRepo.List count execute: %w", err)
	}

	if total == 0 {
		return []domain.PopulatedStudent{}, 0, nil // return an empty slice
	}

	query := `
		SELECT
			s.id, s.first_name, s.last_name, s.email, s.is_active, s.created_at, s.updated_at,
			c.id AS "class.id",
			c.grade AS "class.grade",
			c.letter AS "class.letter"
		FROM students s
		INNER JOIN classes c ON s.class_id = c.id
		ORDER BY s.id ASC
		LIMIT ? OFFSET ?
	`

	students := make([]domain.PopulatedStudent, 0, limit)
	if err := r.db.SelectContext(ctx, &students, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("studentRepo.List select execute: %w", err)
	}

	return students, total, nil
}

func (r *studentRepo) Update(ctx context.Context, s *domain.Student) error {
	query := `
		UPDATE students SET
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			class_id = :class_id,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("studentRepo.Update execute: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("studentRepo.Update rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("studentRepo.Update: %w", domain.ErrNotFound)
	}

	return nil
}
