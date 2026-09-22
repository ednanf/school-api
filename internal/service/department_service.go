package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type departmentService struct {
	repo  domain.DepartmentRepository
	caser cases.Caser
}

func NewDepartmentService(repo domain.DepartmentRepository) domain.DepartmentService {
	return &departmentService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *departmentService) Create(ctx context.Context, d *domain.Department) error {
	// Normalize fields
	d.Normalize(s.caser)

	// Set timestamps
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now

	// Save to database via repository
	if err := s.repo.Create(ctx, d); err != nil {
		return fmt.Errorf("departmentService.Create: %w", err)
	}

	return nil
}

func (s *departmentService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("departmentService.Delete: %w", err)
	}
	return nil
}

func (s *departmentService) GetById(ctx context.Context, id int) (*domain.Department, error) {
	dept, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("departmentService.GetById: %w", err)
	}

	return dept, nil
}

func (s *departmentService) List(ctx context.Context, page, limit int) ([]domain.Department, int, int, int, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	departments, totalItems, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("departmentService.List: %w", err)
	}

	return departments, totalItems, page, limit, nil
}

func (s *departmentService) Update(ctx context.Context, id int, input domain.PatchDepartmentInput) (*domain.Department, error) {
	// Fetch current record
	department, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("departmentService.Update fetch: %w", err)
	}

	// Overwrite raw fields if provided in payload
	if input.Name != nil {
		department.Name = *input.Name
	}

	if input.Description != nil {
		department.Description = *input.Description
	}

	// Normalize data
	department.Normalize(s.caser)

	// Update timestamp
	department.UpdatedAt = time.Now().UTC()

	// Save changes back to DB
	if err := s.repo.Update(ctx, department); err != nil {
		return nil, fmt.Errorf("departmentService.Update execute: %w", err)
	}

	return department, nil
}
