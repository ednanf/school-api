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
