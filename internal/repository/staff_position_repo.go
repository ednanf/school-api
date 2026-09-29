package repository

import (
	"context"

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
