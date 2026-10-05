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
	e.CreatedAt = now
	e.UpdatedAt = now

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

func (s *staffService) GetByID(ctx context.Context, id int) (*domain.PopulatedStaff, error) {
	staff, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("staffService.GetByID: %w", err)
	}
	return staff, nil
}

func (s *staffService) List(ctx context.Context, limit, offset int) (staff []domain.PopulatedStaff, totalItems int, err error) {
	// Enforce defensive pagination boundaries
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	// Delegate fetching and total count to repo layer
	staff, totalItems, err = s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("staffService.List: %w", err)
	}

	return staff, totalItems, nil
}

func (s *staffService) Update(ctx context.Context, id int, input domain.PatchStaffInput) (*domain.Staff, error) {
	// Fetch current record
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("staffService.Update fetch: %w", err)
	}

	// Map the existing record to domain.Staff
	employee := domain.Staff{
		ID:             existing.ID,
		PositionID:     existing.Position.ID,
		FirstName:      existing.FirstName,
		LastName:       existing.LastName,
		Email:          existing.Email,
		HireDateString: existing.HireDate.Format("2006-01-02"), // Preserve original date in string format if not provided in patch
		HireDate:       existing.HireDate,
		IsActive:       existing.IsActive,
		CreatedAt:      existing.CreatedAt,
		UpdatedAt:      existing.UpdatedAt,
	}

	// Apply non-nil updates
	if input.PositionID != nil {
		employee.PositionID = *input.PositionID
	}
	if input.FirstName != nil {
		employee.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		employee.LastName = *input.LastName
	}
	if input.Email != nil {
		employee.Email = *input.Email
	}
	if input.HireDateString != nil {
		// Parse date coming from the JSON payload
		parsedDate, err := time.Parse("2006-01-02", *input.HireDateString)
		if err != nil {
			return nil, fmt.Errorf("staffService.Update: invalid hire_date format: %w", err)
		}

		// Set it in the actual field used in the database
		employee.HireDate = parsedDate

		// Necessary so the date sent throught the payload appears back in the success response (will simply put the same date that came from the payload in the JSON to be sent back)
		employee.HireDateString = *input.HireDateString
	}
	if input.IsActive != nil {
		employee.IsActive = *input.IsActive
	}

	// Normalize and timestamp
	employee.Normalize(s.caser)
	employee.UpdatedAt = time.Now().UTC()

	// Persist merged entity via repo layer
	if err := s.repo.Update(ctx, &employee); err != nil {
		return nil, fmt.Errorf("staffService.Update: %w", err)
	}

	return &employee, nil
}
