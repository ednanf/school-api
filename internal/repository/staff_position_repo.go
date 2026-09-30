package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// staffPositionRepo holds the db connection and methods
type staffPositionRepo struct {
	db *sqlx.DB
}

// NewStaffPositionRepository receives a pointer ot the database connection pool and returns a domain.NewStaffPositionRepository, guaranteeing implementation of all methods
func NewStaffPositionRepository(db *sqlx.DB) domain.StaffPositionRepository {
	return &staffPositionRepo{db: db}
}

func (r *staffPositionRepo) Create(ctx context.Context, s *domain.StaffPosition) error {
	query := `
		INSERT INTO staff_positions (department_id, title, description, is_active, created_at, updated_at)
		VALUES (:department_id, :title, :description, :is_active, :created_at, :updated_at)
	`

	// s is already normalized and timestamped in the service layer
	result, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("staffPositionRepo.Create execute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("staffPositionRepo.Create last insert id: %w", err)
	}

	s.ID = int(id)
	return nil
}

func (r *staffPositionRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM staff_positions WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("staffPositionRepo.Delete execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("staffPositionRepo.Delete rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("staffPositionRepo.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *staffPositionRepo) GetByID(ctx context.Context, id int) (*domain.PopulatedStaffPosition, error) {
	query := `
		SELECT
			s.id, s.title, s.description, s.is_active, s.created_at, s.updated_at,
			d.id AS "department.id",
			d.name AS "department.name",
			d.description AS "department.description"
		FROM staff_positions s
		INNER JOIN departments d ON s.department_id = d.id
		WHERE s.id = ?
	`

	var position domain.PopulatedStaffPosition
	if err := r.db.GetContext(ctx, &position, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("staffPositionRepo.GetByID: %w", err)
		}
		return nil, fmt.Errorf("staffPositionRepo.GetByID execute: %w", err)
	}

	return &position, nil
}

func (r *staffPositionRepo) List(ctx context.Context, limit, offset int) (positions []domain.PopulatedStaffPosition, totalItems int, err error) {
	countQuery := "SELECT COUNT(*) FROM staff_positions"
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("staffPositionRepo.List count execute: %w", err)
	}

	if total == 0 {
		return []domain.PopulatedStaffPosition{}, 0, nil // return an empty slice
	}

	query := `
		SELECT
			s.id, s.title, s.description, s.is_active, s.created_at, s.updated_at,
			d.id AS "department.id",
			d.name AS "department.name",
			d.description AS "department.description"
		FROM staff_positions s
		INNER JOIN departments d ON s.department_id = d.id
		ORDER BY s.id ASC
		LIMIT ? OFFSET ?
	`

	positions = make([]domain.PopulatedStaffPosition, 0, limit)
	if err := r.db.SelectContext(ctx, &positions, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("staffPositionRepo.List select execute: %w", err)
	}

	return positions, total, nil
}

func (r *staffPositionRepo) Update(ctx context.Context, id int) error {
	return nil
}
