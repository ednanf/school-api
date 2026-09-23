package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type teacherAssignmentService struct {
	repo  domain.TeacherAssignmentRepository
	caser cases.Caser
}

func NewTeacherAssignmentRepository(repo domain.TeacherAssignmentRepository) domain.TeacherAssignmentService {
	return &teacherAssignmentService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *teacherAssignmentService) Create(ctx context.Context, t *domain.TeacherAssignment) error {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	if err := s.repo.Create(ctx, t); err != nil {
		return fmt.Errorf("teacherAssignmentService.Create: %w", err)
	}

	return nil
}

func (s *teacherAssignmentService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("teacherAssignmentService.Delete: %w", err)
	}
	return nil
}

func (s *teacherAssignmentService) GetById(ctx context.Context, id int) (*domain.PopulatedTeacherAssignment, error) {
	teacherAssignment, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("teacherAssignmentService.GetById: %w", err)
	}

	return teacherAssignment, nil
}

func (s *teacherAssignmentService) List(ctx context.Context, page, limit int) (items []domain.PopulatedTeacherAssignment, totalItems, pageNum, limitNum int, err error) {
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

	teacherAssignments, totalItems, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("teacherAssignmentsService.List: %w", err)
	}

	return teacherAssignments, totalItems, page, limit, nil
}

func (s *teacherAssignmentService) Update(ctx context.Context, id int, input domain.PatchTeacherAssignmentInput) (*domain.PopulatedTeacherAssignment, error) {
	// Save changes back to db and return hydrated populated result
	assignment, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, fmt.Errorf("teacherAssignmentService.Update execute: %w", err)
	}

	return assignment, nil
}
