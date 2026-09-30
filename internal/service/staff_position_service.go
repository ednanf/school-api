package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type staffPositionService struct {
	repo  domain.StaffPositionRepository
	caser cases.Caser
}

func NewStaffPositionService(repo domain.StaffPositionRepository) domain.StaffPositionService {
	return &staffPositionService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *staffPositionService) Create(ctx context.Context, staffPosition *domain.StaffPosition) error {
	now := time.Now().UTC()

	// Normalize text fields and assign timestamps
	staffPosition.Normalize(s.caser)
	staffPosition.CreatedAt = now
	staffPosition.UpdatedAt = now

	// Delegate persistence to repo with layer error wrapping
	if err := s.repo.Create(ctx, staffPosition); err != nil {
		return fmt.Errorf("staffPositionService.Create:  %w", err)
	}

	return nil
}

func (s *staffPositionService) Delete(ctx context.Context, id int) error {
	// Delegate persistence to repo layer with error wrapping
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("staffPositionService.Delete: %w", err)
	}

	return nil
}

func (s *staffPositionService) GetByID(ctx context.Context, id int) (*domain.PopulatedStaffPosition, error) {
	position, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("staffPositionService.GetByID: %w", err)
	}

	return position, nil
}

func (s *staffPositionService) List(ctx context.Context, limit, offset int) (positions []domain.PopulatedStaffPosition, totalItems int, err error) {
	// Enforce defensive pagination boundaries
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	// Delegate fetching and total count to repository layer
	positions, totalItems, err = s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("staffPositions.List: %w", err)
	}

	return positions, totalItems, nil
}

func (s *staffPositionService) Update(ctx context.Context, id int) (*domain.StaffPosition, error) {
	return nil, nil
}
