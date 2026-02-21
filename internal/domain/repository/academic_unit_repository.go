package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// AcademicUnitRepository defines persistence operations for AcademicUnit
type AcademicUnitRepository interface {
	Create(ctx context.Context, unit *entities.AcademicUnit) error
	FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error)
	FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error)
	FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error)
	Update(ctx context.Context, unit *entities.AcademicUnit) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error)
	ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error)
}
