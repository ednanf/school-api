package repository

import (
	"context"
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
	return nil
}

func (r *staffPositionRepo) GetByID(ctx context.Context, id int) (*domain.PopulatedStaffPosition, error) {
	return nil, nil
}

func (r *staffPositionRepo) List(ctx context.Context, limit, offset int) (positions []domain.PopulatedStaffPosition, totalItems int, err error) {
	return nil, 0, nil
}

func (r *staffPositionRepo) Update(ctx context.Context, id int) error {
	return nil
}
