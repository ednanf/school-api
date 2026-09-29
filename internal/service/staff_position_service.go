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
	return nil
}

func (s *staffPositionService) GetByID(ctx context.Context, id int) (*domain.PopulatedStaffPosition, error) {
	return nil, nil
}

func (s *staffPositionService) List(ctx context.Context, limit, offset int) (positions []domain.PopulatedStaffPosition, totalItems int, err error) {
	return nil, 0, nil
}

func (s *staffPositionService) Update(ctx context.Context, id int) (*domain.StaffPosition, error) {
	return nil, nil
}
