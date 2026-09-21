package repository

import (
	"context"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type departmentRepo struct {
	db *sqlx.DB
}

func NewDepartmentRepository(db *sqlx.DB) domain.DepartmentRepository {
	return &departmentRepo{db: db}
}

func (r *departmentRepo) Create(ctx context.Context, d domain.Department) error {
	return nil
}

func (r *departmentRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *departmentRepo) GetById(ctx context.Context, id int) (*domain.Department, error) {
	return nil, nil
}

func (r *departmentRepo) List(ctx context.Context, limit int, offset int) ([]domain.Department, int, error) {
	return nil, 0, nil
}

func (r *departmentRepo) Update(ctx context.Context, id int, input domain.PatchDepartmentInput) (*domain.Department, error) {
	return nil, nil
}
