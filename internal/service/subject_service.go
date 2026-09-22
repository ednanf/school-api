package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type subjectService struct {
	repo  domain.SubjectRepository
	caser cases.Caser
}

func NewSubjectService(repo domain.SubjectRepository) domain.SubjectService {
	return &subjectService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *subjectService) Create(ctx context.Context, sub *domain.Subject) error {
	// Normalize
	sub.Normalize(s.caser)

	// Set timestamps
	now := time.Now().UTC()
	sub.CreatedAt = now
	sub.UpdatedAt = now

	// Save to database via repository
	if err := s.repo.Create(ctx, sub); err != nil {
		return fmt.Errorf("subjectService.Create: %w", err)
	}

	return nil
}

func (s *subjectService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("subjectService.Delete: %w", err)
	}

	return nil
}

func (s *subjectService) GetById(ctx context.Context, id int) (*domain.Subject, error) {
	// Execute the repository step
	subject, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("subjectService.GetById: %w", err)
	}

	return subject, nil
}

func (s *subjectService) List(ctx context.Context, page, limit int) (items []domain.Subject, totalItems, pageNum, limitNum int, err error) {
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

	subjects, totalItems, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("subjectService.List: %w", err)
	}

	return subjects, totalItems, page, limit, nil
}

func (s *subjectService) Update(ctx context.Context, id int, input domain.PatchSubjectInput) (*domain.Subject, error) {
	// Fetch current record
	subject, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("subjectService.Update fetch: %w", err)
	}

	// Overwrite raw fields if provided in the payload
	if input.Name != nil {
		subject.Name = *input.Name
	}

	if input.IsActive != nil {
		subject.IsActive = *input.IsActive
	}

	// Normalize data
	subject.Normalize(s.caser)

	// Update timestamp
	subject.UpdatedAt = time.Now().UTC()

	// Save changes back to db
	if err := s.repo.Update(ctx, subject); err != nil {
		return nil, fmt.Errorf("subjectService.Update execute: %w", err)
	}

	return subject, nil
}
