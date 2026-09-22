package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type departmentRepo struct {
	db *sqlx.DB
}

func NewDepartmentRepository(db *sqlx.DB) domain.DepartmentRepository {
	return &departmentRepo{db: db}
}

func (r *departmentRepo) Create(ctx context.Context, d *domain.Department) error {
	query := `
		INSERT INTO departments (name, description, created_at, updated_at)
		VALUES (:name, :description, :created_at, :updated_at)
	`

	result, err := r.db.NamedExecContext(ctx, query, d)
	if err != nil {
		return fmt.Errorf("departmentRepo.Create execute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("departmentRepo.Create last insert id: %w", err)
	}

	d.ID = int(id)
	return nil
}

func (r *departmentRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM departments WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("departmentRepo.Delete execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("departmentRepo.Delete rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *departmentRepo) GetById(ctx context.Context, id int) (*domain.Department, error) {
	var d domain.Department

	query := "SELECT id, name, description, created_at, updated_at FROM departments WHERE id = ?"

	if err := r.db.GetContext(ctx, &d, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("departmentRepo.GetById execute: %w", err)
	}

	return &d, nil
}

func (r *departmentRepo) List(ctx context.Context, limit int, offset int) ([]domain.Department, int, error) {
	// departments := make([]domain.Department, 0)

	// var totalItems int
	// countQuery := "SELECT COUNT(*) FROM departments"

	return nil, 0, nil
}

func (r *departmentRepo) Update(ctx context.Context, id int, input domain.PatchDepartmentInput) (*domain.Department, error) {
	return nil, nil
}
