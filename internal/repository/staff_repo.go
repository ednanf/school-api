package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepository(db *sqlx.DB) domain.StaffRepository {
	return &staffRepo{db: db}
}

func (r *staffRepo) Create(ctx context.Context, e *domain.Staff) error {
	query := `
		INSERT INTO staff (position_id, first_name, last_name, email, hire_date, is_active, created_at, updated_at)
		VALUES (:position_id, :first_name, :last_name, :email, :hire_date, :is_active, :created_at, :updated_at)
	`

	// s is normalized and timestampd
	result, err := r.db.NamedExecContext(ctx, query, e)
	if err != nil {
		return fmt.Errorf("staffRepo.Create execute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("staffPositionRepo.Create last insert id: %w", err)
	}

	e.ID = int(id)

	return nil
}

func (r *staffRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM staff WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("staffRepo.Delete execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("staffRepo.Delete rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("staffRepo.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *staffRepo) GetByID(ctx context.Context, id int) (*domain.PopulatedStaff, error) {
	query := `
		SELECT
			-- Staff
			s.id,
	 		s.first_name,
    		s.last_name,
      		s.email,
        	s.hire_date,
         	s.is_active,
          	s.created_at,
           	s.updated_at,

            -- Staff Position Summary
            sp.id AS "position.id",
            sp.title AS "position.title",
            sp.description AS "position.description",

            -- Department Summary
            d.id AS "position.department.id",
            d.name AS "position.department.name",
            d.description AS "position.department.description"
        FROM
        	staff s
       	LEFT JOIN
            staff_positions sp ON s.position_id = sp.id
        LEFT JOIN
            departments d ON sp.department_id = d.id
        WHERE
        	s.id = ?
	`

	var employee domain.PopulatedStaff
	if err := r.db.GetContext(ctx, &employee, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("staffRepo.GetByID: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("staffRepo.GetByID execute: %w", err)
	}

	return &employee, nil
}

func (r *staffRepo) List(ctx context.Context, limit, offset int) (staff []domain.PopulatedStaff, totalItems int, err error) {
	countQuery := "SELECT COUNT(*) FROM staff"
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("staffPositionRepo.List count execute: %w", err)
	}

	if total == 0 {
		return []domain.PopulatedStaff{}, 0, nil // return an empty slice
	}

	query := `
		SELECT
			-- Staff
			s.id,
			s.first_name,
    		s.last_name,
      		s.email,
        	s.hire_date,
         	s.is_active,
          	s.created_at,
           	s.updated_at,

            -- Staff Position Summary
            sp.id AS "position.id",
            sp.title AS "position.title",
            sp.description AS "position.description",

            -- Department Summary
            d.id AS "position.department.id",
            d.name AS "position.department.name",
            d.description AS "position.department.description"
        FROM
        	staff s
       	LEFT JOIN
            staff_positions sp ON s.position_id = sp.id
        LEFT JOIN
            departments d ON sp.department_id = d.id
        ORDER BY
        	s.id ASC
        LIMIT ?
        OFFSET ?
	`

	staff = make([]domain.PopulatedStaff, 0, limit)
	if err := r.db.SelectContext(ctx, &staff, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("staffRepo.List select execute: %w", err)
	}

	return staff, total, nil
}

func (r *staffRepo) Update(ctx context.Context, e *domain.Staff) error {
	query := `
		UPDATE staff SET
			position_id = :position_id,
			first_name = :first_name,
		 	last_name = :last_name,
			email = :email,
		 	hire_date = :hire_date,
			is_active = :is_active,
			created_at = :created_at,
		 	updated_at = :updated_at
		WHERE
			id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, e)
	if err != nil {
		return fmt.Errorf("staffRepo.Update execute: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("staffRepo.Update rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("staffRepo.Update: %w", domain.ErrNotFound)
	}

	return nil
}
