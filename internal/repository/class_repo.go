package repository

import (
	"context"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type classRepo struct {
	db *sqlx.DB
}

// Constructor
func NewClassRepository(db *sqlx.DB) domain.ClassRepository {
	return &classRepo{db: db}
}

func (r *classRepo) Create(ctx context.Context, class *domain.Class) error {
	return nil
}

func (r *classRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *classRepo) GetById(ctx context.Context, id int) (*domain.Class, error) {
	return nil, nil
}

func (r *classRepo) List(ctx context.Context, limit int, offset int) ([]domain.Class, error) {
	return nil, nil
}

func (r *classRepo) Update(ctx context.Context, id int, input domain.PatchClassInput) (*domain.Class, error) {
	return &domain.Class{}, nil
}
