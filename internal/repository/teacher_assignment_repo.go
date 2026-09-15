package repository

import (
	"context"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// taRepo stores the db connection and the repository methods attached to it
type taRepo struct {
	db *sqlx.DB
}

// NewTeacherAssignmentRepository receives a pointer to the database conneciton pool and returns a domain.TeacherAssignmentRepository
func NewTeacherAssignmentRepository(db *sqlx.DB) domain.TeacherAssignmentRepository {
	return &taRepo{db: db}
}

func (r *taRepo) Create(ctx context.Context, t *domain.TeacherAssignment) error {
	return nil
}

func (r *taRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *taRepo) GetById(ctx context.Context, id int) (*domain.TeacherAssignment, error) {
	return nil, nil
}

// List takes a context, limit and offset and returns a slice, a total and errors
func (r *taRepo) List(ctx context.Context, limit, offset int) ([]domain.TeacherAssignment, int, error) {
	return nil, 0, nil
}

func (r *taRepo) Update(ctx context.Context, id int, input domain.PatchTeacherAssignmentInput) (*domain.TeacherAssignment, error) {
	return nil, nil
}
