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

// studentRepo stores the db connnection and the repository methods attached to it
type studentRepo struct {
	db *sqlx.DB
}

// NewStudentRepository receives a pointer to the database connection pool and returns a domain.StudentRepository, guaranteeing studentRepo implements all required methods
func NewStudentRepository(db *sqlx.DB) domain.StudentRepository {
	return &studentRepo{db: db}
}

// BatchCreate receives a context and a slice of type Student. It accepts 1000 entries at most. It returns a slice containing all created students, the total amount of created entries and errors
func (r *studentRepo) BatchCreate(ctx context.Context, students []domain.Student) ([]domain.Student, int, error) {
	if len(students) == 0 {
		return students, 0, nil
	}

	// Initiate a db transaction with the context
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("studentRepo.BatchCreate begin tx: %w", err)
	}

	// Should any error occur, roll back the database
	defer tx.Rollback()

	query := `
        INSERT INTO students (first_name, last_name, email, class_id, created_at, updated_at)
        VALUES (:first_name, :last_name, :email, :class_id, :created_at, :updated_at)
    `
	caser := cases.Title(language.English)
	now := time.Now().UTC()
	total := 0

	// Loop the student slice argument
	for i := range students {
		// Add timestamps
		students[i].Normalize(caser)
		students[i].CreatedAt = now
		students[i].UpdatedAt = now

		// Execute the database operation in the ongoing transaction
		res, err := tx.NamedExecContext(ctx, query, students[i])
		if err != nil {
			return nil, 0, fmt.Errorf("studentRepo.BatchCreate insert: %w", err)
		}

		// Retrieve the newly inserted entry's id to be able to send a response
		id, err := res.LastInsertId()
		if err != nil {
			return nil, 0, fmt.Errorf("studentRepo.BatchCreate last insert id: %w", err)
		}

		// Assign the received id to the entry in order to show in the response
		students[i].ID = int(id)

		total++
	}

	// Commit the database changes
	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("studentRepo.BatchCreate commit: %w", err)
	}

	return students, total, nil
}

// BatchDelete receives a context and a slice of ids of type int. It accepts a maximum of 100 ids. It returns a total and errors
func (r *studentRepo) BatchDelete(ctx context.Context, ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Expand the slice into positional placeholders: WHERE id IN (?, ?, ...)
	query, args, err := sqlx.In("DELETE FROM students WHERE id IN (?)", ids)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BatchDelete query build: %w", err)
	}

	// Rebind query to match MariaDB driver syntax
	query = r.db.Rebind(query)

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BatchDelete execute: %w", err)
	}

	// Return total number of deleted rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BatchDelete rows affected: %w", err)
	}

	return rowsAffected, nil
}

// BatchUpdate receives a slice of batch inputs and updates each student record in a single transaction.
func (r *studentRepo) BatchUpdate(ctx context.Context, updates []domain.BatchUpdateStudentItem) ([]domain.Student, int, error) {
	if len(updates) == 0 {
		return []domain.Student{}, 0, nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("studentRepo.BatchUpdate begin tx: %w", err)
	}
	defer tx.Rollback()

	caser := cases.Title(language.English)
	now := time.Now().UTC()
	updatedStudents := make([]domain.Student, 0, len(updates))

	query := `
		UPDATE students SET
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			class_id = :class_id,
			updated_at = :updated_at
		WHERE id = :id
	`

	total := 0

	for _, u := range updates {
		// Fetch current record inside the active transaction
		var student domain.Student
		fetchQuery := "SELECT * FROM students WHERE id = ?"
		if err := tx.GetContext(ctx, &student, fetchQuery, u.ID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, 0, fmt.Errorf("studentRepo.BatchUpdate student ID %d not found: %w", u.ID, sql.ErrNoRows)
			}
			return nil, 0, fmt.Errorf("studentRepo.BatchUpdate fetch ID %d: %w", u.ID, err)
		}

		// Apply non-nil updates (promoted directly from embedded PatchStudentInput)
		if u.FirstName != nil {
			student.FirstName = *u.FirstName
		}
		if u.LastName != nil {
			student.LastName = *u.LastName
		}
		if u.Email != nil {
			student.Email = *u.Email
		}
		if u.ClassID != nil {
			student.ClassID = *u.ClassID
		}

		// Normalize fields and set updated timestamp
		student.Normalize(caser)
		student.UpdatedAt = now

		// Execute update query
		_, err = tx.NamedExecContext(ctx, query, student)
		if err != nil {
			return nil, 0, fmt.Errorf("studentRepo.BatchUpdate execute ID %d: %w", u.ID, err)
		}

		updatedStudents = append(updatedStudents, student)
		total++
	}

	// Commit transaction changes
	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("studentRepo.BatchUpdate commit: %w", err)
	}

	return updatedStudents, total, nil
}

// BulkUpdateClass reassigns a slice of student IDs to a new class_id in a single execution.
func (r *studentRepo) BulkUpdateClass(ctx context.Context, ids []int, classID int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// 1. Verify target class exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM classes WHERE id = ?)"
	if err := r.db.GetContext(ctx, &exists, checkQuery, classID); err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdateClass check class: %w", err)
	}
	if !exists {
		return 0, domain.ErrClassNotFound
	}

	// 2. Perform bulk update
	rawQuery := "UPDATE students SET class_id = ?, updated_at = ? WHERE id IN (?)"
	now := time.Now().UTC()

	query, args, err := sqlx.In(rawQuery, classID, now, ids)
	if err != nil {
		return 0, fmt.Errorf("studentRepo.BulkUpdateClass query build: %w", err)
	}

	query = r.db.Rebind(query)

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
        INSERT INTO students (first_name, last_name, email, class_id, created_at, updated_at)
        VALUES (:first_name, :last_name, :email, :class_id, :created_at, :updated_at)
    `

	// Initiate a caser to normalize capitalization
	caser := cases.Title(language.English)

	// Normalize fields explicitly (Trim + Case handling)
	s.Normalize(caser)

	// Add timestamp
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now

	// Execute the db operation with `NamedExecContext` to match the named placeholders
	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("studentRepo.Create execute: %w", err)
	}

	// Grab newly inserted entry's id to be able to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("studentRepo.Create last insert id: %w", err)
	}

	// Assign the received id to the entry in order to show in the response
	s.ID = int(id)

	return nil
}

func (r *studentRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM students WHERE id = ?"

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("studentRepo.Delete execute: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("studentRepo.Delete rows affected: %w", err)
	}

	// If 0 rows were affected, the ID did not exist in the db
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *studentRepo) GetByID(ctx context.Context, id int) (*domain.Student, error) {
	var s domain.Student

	query := "SELECT * FROM students WHERE id = ?"

	// Execute query and assign it to the variable `s` if successful
	err := r.db.GetContext(ctx, &s, query, id)
	if err != nil {
		// If the id is not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// Other errors
		return nil, fmt.Errorf("studentRepo.GetByID execute: %w", err)
	}

	return &s, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *studentRepo) List(ctx context.Context, limit int, offset int) ([]domain.Student, int, error) {
	// Make an empty slice to hold students
	students := make([]domain.Student, 0)

	// Get the total count across the entire table
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM students"
	if err := r.db.GetContext(ctx, &totalItems, countQuery); err != nil {
		return nil, 0, fmt.Errorf("studentRepo.List count: %w", err)
	}

	// Get all columns from the table students, ordered by their ID, and limited to a certain number
	query := "SELECT * FROM students ORDER BY id LIMIT ? OFFSET ?"

	// Execute the db operation
	err := r.db.SelectContext(ctx, &students, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("studentRepo.List fetch: %w", err)
	}

	// Return the results to be used
	return students, totalItems, nil
}

func (r *studentRepo) Update(ctx context.Context, id int, input domain.PatchStudentInput) (*domain.Student, error) {
	// Fetch current record from DB
	student, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("studentRepo.Update fetch: %w", err)
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

	// Normalize the entity as a whole
	caser := cases.Title(language.English)
	student.Normalize(caser)

	// Apply timestamp
	student.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE students SET
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			class_id = :class_id,
			updated_at = :updated_at
		WHERE id = :id
	`

	// Execute static SQL query using sqlx named placeholders
	_, err = r.db.NamedExecContext(ctx, query, student)
	if err != nil {
		return nil, fmt.Errorf("studentRepo.Update execute: %w", err)
	}

	return student, nil
}
