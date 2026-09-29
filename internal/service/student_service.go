package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ednanf/school-api/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type studentService struct {
	repo  domain.StudentRepository
	caser cases.Caser
}

func NewStudentService(repo domain.StudentRepository) domain.StudentService {
	return &studentService{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *studentService) BulkCreate(ctx context.Context, students []domain.Student) ([]domain.Student, int, error) {
	if len(students) == 0 {
		return students, 0, nil
	}

	now := time.Now().UTC()

	// Business logic & metadata preparation
	for i := range students {
		students[i].Normalize(s.caser)
		students[i].CreatedAt = now
		students[i].UpdatedAt = now
	}

	// Database persistence with explicit error wrapping
	createdStudents, total, err := s.repo.BulkCreate(ctx, students)
	if err != nil {
		return nil, 0, fmt.Errorf("studentService.BulkCreate: %w", err)
	}

	return createdStudents, total, nil
}

func (s *studentService) BulkDelete(ctx context.Context, ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	rowsAffected, err := s.repo.BulkDelete(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("studentService.BulkDelete: %w", err)
	}

	return rowsAffected, nil
}

func (s *studentService) BulkUpdate(ctx context.Context, updates []domain.BulkUpdateStudentItem) ([]domain.Student, int, error) {
	if len(updates) == 0 {
		return []domain.Student{}, 0, nil
	}

	now := time.Now().UTC()
	studentsToUpdate := make([]domain.Student, 0, len(updates))

	// Loop through the updates slice and append the updated entry to `studentsToUpdate`
	for _, u := range updates {
		// Fetch current record from repository
		existing, err := s.repo.GetByID(ctx, u.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("studentService.BulkUpdate fetch ID %d: %w", u.ID, err)
		}
		if existing == nil {
			return nil, 0, fmt.Errorf("studentService.BulkUpdate student ID %d not found: %w", u.ID, domain.ErrNotFound)
		}

		// Map existing record to domain model
		student := domain.Student{
			ID:        existing.ID,
			FirstName: existing.FirstName,
			LastName:  existing.LastName,
			Email:     existing.Email,
			ClassID:   existing.Class.ID,
			IsActive:  existing.IsActive,
			CreatedAt: existing.CreatedAt,
			UpdatedAt: existing.UpdatedAt,
		}

		// Apply non-nil patch field updates
		if u.FirstName != nil {
			student.FirstName = *u.FirstName
		}
		if u.LastName != nil {
			student.LastName = *u.LastName
		}
		if u.Email != nil {
			student.Email = *u.Email
		}
		if u.ClassID != nil {
			student.ClassID = *u.ClassID
		}
		if u.IsActive != nil {
			student.IsActive = *u.IsActive
		}

		// Normalize text and update timestamp
		student.Normalize(s.caser)
		student.UpdatedAt = now

		studentsToUpdate = append(studentsToUpdate, student)
	}

	// Pass fully merged & normalized entities to repo for single-transaction update
	_, err := s.repo.BulkUpdate(ctx, studentsToUpdate)
	if err != nil {
		return nil, 0, fmt.Errorf("studentService.BulkUpdate: %w", err)
	}

	return studentsToUpdate, len(studentsToUpdate), nil
}

func (s *studentService) BulkUpdateClass(ctx context.Context, ids []int, classID int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// 1. If you have a classRepo injected into studentService, check class existence here:
	// exists, err := s.classRepo.Exists(ctx, classID)
	// if err != nil {
	// 	return 0, fmt.Errorf("studentService.BulkUpdateClass check class: %w", err)
	// }
	// if !exists {
	// 	return 0, domain.ErrNotFound
	// }

	// Generate UTC timestamp in service layer
	now := time.Now().UTC()

	// Pass IDs, target class, and timestamp to repository
	rowsAffected, err := s.repo.BulkUpdateClass(ctx, ids, classID, now)
	if err != nil {
		return 0, fmt.Errorf("studentService.BulkUpdateClass: %w", err)
	}

	return rowsAffected, nil
}

func (s *studentService) Create(ctx context.Context, student *domain.Student) error {
	now := time.Now().UTC()

	// Normalize text fields and assign timestamps
	student.Normalize(s.caser)
	student.CreatedAt = now
	student.UpdatedAt = now

	// Delegate persistence to repo with layer error wrapping
	if err := s.repo.Create(ctx, student); err != nil {
		return fmt.Errorf("studentService.Create: %w", err)
	}

	return nil
}

func (s *studentService) Delete(ctx context.Context, id int) error {
	// Delegate persistence to repo with layer error wrapping
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("studentService.Delete: %w", err)
	}

	return nil
}

func (s *studentService) GetByID(ctx context.Context, id int) (*domain.PopulatedStudent, error) {
	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("studentService.GetByID: %w", err)
	}

	return student, nil
}

func (s *studentService) List(ctx context.Context, limit, offset int) ([]domain.PopulatedStudent, int, error) {
	// Enforce defensive pagination boundaries in the business layer
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	// Delegate fetching and total count to repository
	students, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("studentService.List: %w", err)
	}

	return students, total, nil
}

func (s *studentService) Update(ctx context.Context, id int, input domain.PatchStudentInput) (*domain.Student, error) {
	// Fetch current record to build full entity state
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("studentService.Update fetch: %w", err)
	}

	// Map existing PopulatedStudent record to domain.Student
	student := domain.Student{
		ID:        existing.ID,
		FirstName: existing.FirstName,
		LastName:  existing.LastName,
		Email:     existing.Email,
		ClassID:   existing.Class.ID,
		IsActive:  existing.IsActive,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}

	// Apply non-nil patch updates
	if input.FirstName != nil {
		student.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		student.LastName = *input.LastName
	}
	if input.Email != nil {
		student.Email = *input.Email
	}
	if input.ClassID != nil {
		student.ClassID = *input.ClassID
	}
	if input.IsActive != nil {
		student.IsActive = *input.IsActive
	}

	// Normalize string fields and update timestamp
	student.Normalize(s.caser)
	student.UpdatedAt = time.Now().UTC()

	// Persist merged entity via repository
	if err := s.repo.Update(ctx, &student); err != nil {
		return nil, fmt.Errorf("studentService.Update: %w", err)
	}

	return &student, nil
}
