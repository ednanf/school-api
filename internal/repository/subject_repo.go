package repository

import (
	"context"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

// subjectRepo stores the db connection and has the repository methods attached to it
type subjectRepo struct {
	db *sqlx.DB
}

// NewSubjectRepository receives a pointer to the database connection pool and returns a domain.SubjectRepository, guaranteeing subjectRepo implements all required methods
func NewSubjectRepository(db *sqlx.DB) domain.SubjectRepository {
	return &subjectRepo{db: db}
}

func (r *subjectRepo) Create(ctx context.Context, s *domain.Subject) error {
	return nil
}

func (r *subjectRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *subjectRepo) GetById(ctx context.Context, id int) (*domain.Class, error) {
	return nil, nil
}

func (r *subjectRepo) List(ctx context.Context, limit int, offset int) ([]domain.Subject, int, error) {
	return nil, 0, nil
}

func (r *subjectRepo) Update(ctx context.Context, id int, input domain.PatchSubjectInput) (*domain.Subject, error) {
	return nil, nil
}
