package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
)

// SubjectRepository defines persistence operations for Subject
type SubjectRepository interface {
	Create(ctx context.Context, subject *entities.Subject) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error)
	Update(ctx context.Context, subject *entities.Subject) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.Subject, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
