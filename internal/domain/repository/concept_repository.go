package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// ConceptTypeRepository defines persistence operations for ConceptType
type ConceptTypeRepository interface {
	FindAll(ctx context.Context) ([]*entities.ConceptType, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.ConceptType, error)
	FindByCode(ctx context.Context, code string) (*entities.ConceptType, error)
	Create(ctx context.Context, ct *entities.ConceptType) error
	Update(ctx context.Context, ct *entities.ConceptType) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// ConceptDefinitionRepository defines persistence operations for ConceptDefinition
type ConceptDefinitionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entities.ConceptDefinition, error)
	FindByTypeID(ctx context.Context, typeID uuid.UUID) ([]*entities.ConceptDefinition, error)
	Create(ctx context.Context, def *entities.ConceptDefinition) error
	Update(ctx context.Context, def *entities.ConceptDefinition) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SchoolConceptRepository defines persistence operations for SchoolConcept
type SchoolConceptRepository interface {
	FindBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*entities.SchoolConcept, error)
	BulkCreate(ctx context.Context, concepts []*entities.SchoolConcept) error
	Update(ctx context.Context, concept *entities.SchoolConcept) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.SchoolConcept, error)
}
