package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// GuardianRepository defines persistence operations for GuardianRelation
type GuardianRepository interface {
	Create(ctx context.Context, relation *entities.GuardianRelation) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error)
	FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error)
	FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error)
	Update(ctx context.Context, relation *entities.GuardianRelation) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error)
}
