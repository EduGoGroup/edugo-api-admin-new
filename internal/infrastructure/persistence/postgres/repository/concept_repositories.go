package repository

import (
	"context"
	"errors"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== ConceptType ====================

type postgresConceptTypeRepository struct{ db *gorm.DB }

func NewPostgresConceptTypeRepository(db *gorm.DB) repository.ConceptTypeRepository {
	return &postgresConceptTypeRepository{db: db}
}

func (r *postgresConceptTypeRepository) FindAll(ctx context.Context) ([]*entities.ConceptType, error) {
	var types []*entities.ConceptType
	err := r.db.WithContext(ctx).Where("is_active = true").Order("name").Find(&types).Error
	return types, err
}

func (r *postgresConceptTypeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.ConceptType, error) {
	var ct entities.ConceptType
	if err := r.db.WithContext(ctx).Where("is_active = true").First(&ct, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ct, nil
}

func (r *postgresConceptTypeRepository) FindByCode(ctx context.Context, code string) (*entities.ConceptType, error) {
	var ct entities.ConceptType
	if err := r.db.WithContext(ctx).Where("code = ? AND is_active = true", code).First(&ct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ct, nil
}

func (r *postgresConceptTypeRepository) Create(ctx context.Context, ct *entities.ConceptType) error {
	return r.db.WithContext(ctx).Create(ct).Error
}

func (r *postgresConceptTypeRepository) Update(ctx context.Context, ct *entities.ConceptType) error {
	return r.db.WithContext(ctx).Save(ct).Error
}

func (r *postgresConceptTypeRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.ConceptType{}).Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

// ==================== ConceptDefinition ====================

type postgresConceptDefinitionRepository struct{ db *gorm.DB }

func NewPostgresConceptDefinitionRepository(db *gorm.DB) repository.ConceptDefinitionRepository {
	return &postgresConceptDefinitionRepository{db: db}
}

func (r *postgresConceptDefinitionRepository) FindByTypeID(ctx context.Context, typeID uuid.UUID) ([]*entities.ConceptDefinition, error) {
	var defs []*entities.ConceptDefinition
	err := r.db.WithContext(ctx).Where("concept_type_id = ?", typeID).Order("sort_order, term_key").Find(&defs).Error
	return defs, err
}

func (r *postgresConceptDefinitionRepository) Create(ctx context.Context, def *entities.ConceptDefinition) error {
	return r.db.WithContext(ctx).Create(def).Error
}

func (r *postgresConceptDefinitionRepository) Update(ctx context.Context, def *entities.ConceptDefinition) error {
	return r.db.WithContext(ctx).Save(def).Error
}

func (r *postgresConceptDefinitionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ConceptDefinition{}, "id = ?", id).Error
}

// ==================== SchoolConcept ====================

type postgresSchoolConceptRepository struct{ db *gorm.DB }

func NewPostgresSchoolConceptRepository(db *gorm.DB) repository.SchoolConceptRepository {
	return &postgresSchoolConceptRepository{db: db}
}

func (r *postgresSchoolConceptRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*entities.SchoolConcept, error) {
	var concepts []*entities.SchoolConcept
	err := r.db.WithContext(ctx).Where("school_id = ?", schoolID).Order("category, term_key").Find(&concepts).Error
	return concepts, err
}

func (r *postgresSchoolConceptRepository) BulkCreate(ctx context.Context, concepts []*entities.SchoolConcept) error {
	if len(concepts) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&concepts).Error
}

func (r *postgresSchoolConceptRepository) Update(ctx context.Context, concept *entities.SchoolConcept) error {
	return r.db.WithContext(ctx).Save(concept).Error
}

func (r *postgresSchoolConceptRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.SchoolConcept, error) {
	var sc entities.SchoolConcept
	if err := r.db.WithContext(ctx).First(&sc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sc, nil
}
