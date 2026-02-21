package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// MembershipRepository defines persistence operations for Membership
type MembershipRepository interface {
	Create(ctx context.Context, membership *entities.Membership) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error)
	FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error)
	FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error)
	FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error)
	Update(ctx context.Context, membership *entities.Membership) error
	Delete(ctx context.Context, id uuid.UUID) error
}
