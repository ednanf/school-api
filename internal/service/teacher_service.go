package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type teacherService struct {
	repo  domain.TeacherRepository
	caser cases.Caser
}

func NewTeacherService(repo domain.TeacherRepository) domain.TeacherService {
	return &teacherService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *teacherService) Create(ctx context.Context, t *domain.Teacher) error {
	// Normalize fields
	t.Normalize(s.caser)

	// Set timestamps
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	// Save to database via repository
	if err := s.repo.Create(ctx, t); err != nil {
		return fmt.Errorf("teacherService.Create: %w", err)
	}

	return nil
}

func (s *teacherService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("teacherService.Delete: %w", err)
	}
	return nil
}

func (s *teacherService) GetById(ctx context.Context, id int) (*domain.Teacher, error) {
	teacher, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("teacherService.GetById: %w", err)
	}
	return teacher, nil
}

func (s *teacherService) List(ctx context.Context, page, limit int) (items []domain.Teacher, totalItems, pageNum, limitNum int, err error) {
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

	teachers, totalItems, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("teacherService.List: %w", err)
	}

	return teachers, totalItems, page, limit, nil
}

func (s *teacherService) Update(ctx context.Context, id int, input domain.PatchTeacherInput) (*domain.Teacher, error) {
	// Fetch current record
	teacher, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("teacherService.Update: %w", err)
	}

	// Overwirte raw fields if provided in payload
	if input.FirstName != nil {
		teacher.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		teacher.LastName = *input.LastName
	}

	if input.Email != nil {
		teacher.Email = *input.Email
	}

	if input.IsActive != nil {
		teacher.IsActive = *input.IsActive
	}

	// Normalize data
	teacher.Normalize(s.caser)

	// Update timestamp
	teacher.UpdatedAt = time.Now().UTC()

	// Save changes back to db
	if err := s.repo.Update(ctx, teacher); err != nil {
		return nil, fmt.Errorf("teacherService.Update execute: %w", err)
	}

	return teacher, nil
}
