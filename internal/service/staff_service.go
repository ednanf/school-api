package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type staffService struct {
	repo  domain.StaffRepository
	caser cases.Caser
}

func NewStaffService(repo domain.StaffRepository) domain.StaffService {
	return &staffService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *staffService) Create(ctx context.Context, e *domain.Staff) error {
	now := time.Now().UTC()

	// Parse date string into time.Time
	parsedDate, err := time.Parse("2006-01-02", e.HireDateString)
	if err != nil {
		return fmt.Errorf("staffService.Create: invalid hire_date format: %w", err)
	}
	e.HireDate = parsedDate

	// Normalize and timestamp
	e.Normalize(s.caser)
	e.Created_At = now
	e.Updated_At = now

	// Persist data via repo layer
	if err := s.repo.Create(ctx, e); err != nil {
		return fmt.Errorf("staffService.Create: %w", err)
	}

	return nil
}

func (s *staffService) Delete(ctx context.Context, id int) error {
	// Delegate persistence to repo layer
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("staffService.Delete: %w", err)
	}

	return nil
}

func (s *staffService) GetById(ctx context.Context, id int) (*domain.PopulatedStaff, error) {
	return nil, nil
}

func (s *staffService) List(ctx context.Context, limit, offset int) (employees []domain.PopulatedStaff, totalItems int, err error) {
	return nil, 0, nil
}

func (s *staffService) Update(ctx context.Context, id int, input domain.PatchStaffInput) (*domain.Staff, error) {
	return nil, nil
}
