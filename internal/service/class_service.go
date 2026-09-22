package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type classService struct {
	repo  domain.ClassRepository
	caser cases.Caser
}

func NewClassService(repo domain.ClassRepository) domain.ClassService {
	return &classService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *classService) Create(ctx context.Context, c *domain.Class) error {
	// Normalize fields
	c.Normalize()

	// Set timestamps
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	// Save to db via repository
	if err := s.repo.Create(ctx, c); err != nil {
		return fmt.Errorf("classService.Create: %w", err)
	}

	return nil
}

func (s *classService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("classService.Delete: %w", err)
	}
	return nil
}

func (s *classService) GetById(ctx context.Context, id int) (*domain.Class, error) {
	class, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("classService.GetById: %w", err)
	}

	return class, nil
}

func (s *classService) List(ctx context.Context, page, limit int) (items []domain.Class, totalItems, pageNum, limitNum int, err error) {
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

	classes, totalItems, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("classService.List: %w", err)
	}

	return classes, totalItems, page, limit, nil
}

func (s *classService) ListStudentsByClassId(ctx context.Context, classID int, page int, limit int) (items []domain.Student, totalItems, pageNum, limitNum int, err error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}

	offset := (page - 1) * limit

	students, totalItems, err := s.repo.ListStudentsByClassId(ctx, classID, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("classService.ListStudentsByClassId: %w", err)
	}

	return students, totalItems, page, limit, nil
}

func (s *classService) Update(ctx context.Context, id int, input domain.PatchClassInput) (*domain.Class, error) {
	// Fetch current record
	class, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("classService.Update fetch: %w", err)
	}

	// Overwrite raw fields if provided in payload
	if input.Grade != nil {
		class.Grade = *input.Grade
	}

	if input.Letter != nil {
		class.Letter = *input.Letter
	}

	if input.IsActive != nil {
		class.IsActive = *input.IsActive
	}

	// Normalize data
	class.Normalize()

	// Update timestamp
	class.UpdatedAt = time.Now().UTC()

	// Save changes back to db
	if err := s.repo.Update(ctx, class); err != nil {
		return nil, fmt.Errorf("classService.Update execute: %w", err)
	}

	return class, nil
}
