package repository

import (
	"context"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// teacherRepo stores the db connection and the repository methods attached to it
type teacherRepo struct {
	db *sqlx.DB
}

// NewTeacherRepository receives a pointer to the database connection pool and returns a domain.TeacherRepository, guaranteeing teacherRepo implements all required methods
func NewTeacherRepository(db *sqlx.DB) domain.TeacherRepository {
	return &teacherRepo{db: db}
}

func (r *teacherRepo) Create(ctx context.Context, t *domain.Teacher) error {
	return nil
}

func (r *teacherRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *teacherRepo) GetById(ctx context.Context, id int) (*domain.Teacher, error) {
	return nil, nil
}

func (r *teacherRepo) List(ctx context.Context, limit int, offset int) ([]domain.Teacher, int, error) {
	return nil, 0, nil
}

func (r *teacherRepo) Update(ctx context.Context, id int, input domain.PatchTeacherInput) (*domain.Teacher, error) {
	return nil, nil
}
