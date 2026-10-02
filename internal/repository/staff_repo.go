package repository

import (
	"context"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// TODO: Staff has to populate with StaffPosition, which populates with Departments. Study how to properly populate everything

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepository(db *sqlx.DB) domain.StaffRepository {
	return &staffRepo{db: db}
}

func (r *staffRepo) Create(ctx context.Context, s *domain.StaffPosition) error {
	return nil
}

func (r *staffRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *staffRepo) GetByID(ctx context.Context, id int) (*domain.PopulatedStaff, error) {
	return nil, nil
}

func (r *staffRepo) List(ctx context.Context, page, offset int) (employees []domain.PopulatedStaff, totalItems int, err error) {
	return nil, 0, nil
}

func (r *staffRepo) Update(ctx context.Context, s *domain.Staff) error {
	return nil
}
