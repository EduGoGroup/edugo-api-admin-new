package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// SchoolRepository defines persistence operations for School
type SchoolRepository interface {
	Create(ctx context.Context, school *entities.School) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error)
	FindByCode(ctx context.Context, code string) (*entities.School, error)
	Update(ctx context.Context, school *entities.School) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters ListFilters) ([]*entities.School, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
