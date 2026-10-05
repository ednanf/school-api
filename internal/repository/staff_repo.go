package repository

import (
	"context"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// TODO: Staff has to populate with StaffPosition, which populates with Departments. Study how to properly populate everything

type staffRepo struct {
	db *sqlx.DB
}

// TODO: Start building the repository
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
	return nil, nil
}

func (r *staffRepo) List(ctx context.Context, page, offset int) (employees []domain.PopulatedStaff, totalItems int, err error) {
	return nil, 0, nil
}

func (r *staffRepo) Update(ctx context.Context, e *domain.Staff) error {
	return nil
}
