package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// taRepo stores the db connection and the repository methods attached to it
type taRepo struct {
	db *sqlx.DB
}

// NewTeacherAssignmentRepository receives a pointer to the database connection pool and returns a domain.TeacherAssignmentRepository
func NewTeacherAssignmentRepository(db *sqlx.DB) domain.TeacherAssignmentRepository {
	return &taRepo{db: db}
}

func (r *taRepo) Create(ctx context.Context, t *domain.TeacherAssignment) error {
	query := `
		INSERT INTO teacher_assignments (teacher_id, class_id, subject_id, created_at, updated_at)
		VALUES (:teacher_id, :class_id, :subject_id, :created_at, :updated_at)
	`

	result, err := r.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("taRepo.Create execute: %w", err)
	}

	// Grab newly inserted entry's id to send a response
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("taRepo.Create last inserted id: %w", err)
	}

	// Assign the id to the entry to show in the response
	t.ID = int(id)

	return nil
}

func (r *taRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM teacher_assignments WHERE id = ?"

	// Execute the db operation
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("taRepo.Delete execute: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taRepo.Delete rows affected: %w", err)
	}

	// If 0 rows were affected, the ID did not exist in the db
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *taRepo) GetById(ctx context.Context, id int) (*domain.PopulatedTeacherAssignment, error) {
	var t domain.PopulatedTeacherAssignment

	query := `
		SELECT
			ta.id, ta.created_at, ta.updated_at,
			t.id AS "teacher.id",
			t.first_name AS "teacher.first_name",
			t.last_name AS "teacher.last_name",
			t.email AS "teacher.email",
			c.id AS "class.id",
			c.grade AS "class.grade",
			c.letter AS "class.letter",
			s.id AS "subject.id",
			s.name AS "subject.name"
		FROM teacher_assignments ta
		INNER JOIN teachers t ON ta.teacher_id = t.id
		INNER JOIN classes c ON ta.class_id = c.id
		INNER JOIN subjects s ON ta.subject_id = s.id
		WHERE ta.id = ?
	`

	if err := r.db.GetContext(ctx, &t, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("taRepo.GetByID execute: %w", err)
	}

	return &t, nil
}

// List takes a context, limit and offset and returns a populated slice, a total and errors
func (r *taRepo) List(ctx context.Context, limit, offset int) ([]domain.PopulatedTeacherAssignment, int, error) {
	assignments := make([]domain.PopulatedTeacherAssignment, 0)

	var totalItems int
	countQuery := "SELECT COUNT(*) FROM teacher_assignments"
	if err := r.db.GetContext(ctx, &totalItems, countQuery); err != nil {
		return nil, 0, fmt.Errorf("taRepo.List count: %w", err)
	}

	query := `
		SELECT
			ta.id, ta.created_at, ta.updated_at,
			t.id AS "teacher.id",
			t.first_name AS "teacher.first_name",
			t.last_name AS "teacher.last_name",
			t.email AS "teacher.email",
			c.id AS "class.id",
			c.grade AS "class.grade",
			c.letter AS "class.letter",
			s.id AS "subject.id",
			s.name AS "subject.name"
		FROM teacher_assignments ta
		INNER JOIN teachers t ON ta.teacher_id = t.id
		INNER JOIN classes c ON ta.class_id = c.id
		INNER JOIN subjects s ON ta.subject_id = s.id
		ORDER BY ta.id ASC
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &assignments, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("taRepo.List fetch: %w", err)
	}

	return assignments, totalItems, nil
}

func (r *taRepo) Update(ctx context.Context, id int, input domain.PatchTeacherAssignmentInput) (*domain.PopulatedTeacherAssignment, error) {
	// Make a slice and a map to hold clauses and arguments
	setClauses := make([]string, 0, 4)
	args := make(map[string]any)

	// Append clause and set argument if a value is present
	if input.TeacherID != nil {
		setClauses = append(setClauses, "teacher_id = :teacher_id")
		args["teacher_id"] = *input.TeacherID
	}
	if input.ClassID != nil {
		setClauses = append(setClauses, "class_id = :class_id")
		args["class_id"] = *input.ClassID
	}
	if input.SubjectID != nil {
		setClauses = append(setClauses, "subject_id = :subject_id")
		args["subject_id"] = *input.SubjectID
	}

	// If there are no changes in the input payload
	if len(setClauses) == 0 {
		return r.GetById(ctx, id)
	}

	// Add timestamp clause + argument
	setClauses = append(setClauses, "updated_at = :updated_at")
	args["updated_at"] = time.Now().UTC()
	args["id"] = id

	// Build query dinamically with the values
	query := fmt.Sprintf(`
        UPDATE teacher_assignments
        SET %s
        WHERE id = :id
    `, strings.Join(setClauses, ", "))

	// Execute the operation
	result, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("taRepo.Update execute: %w", err)
	}

	// Determine if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("taRepo.Update rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	// Return the hydrating the nested objects (due to how GetById works)
	return r.GetById(ctx, id)
}
