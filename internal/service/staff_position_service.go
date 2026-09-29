package service

import (
	"context"

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
