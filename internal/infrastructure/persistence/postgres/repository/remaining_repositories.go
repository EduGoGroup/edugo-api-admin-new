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

// ==================== AcademicUnit ====================

type postgresAcademicUnitRepository struct{ db *gorm.DB }

func NewPostgresAcademicUnitRepository(db *gorm.DB) repository.AcademicUnitRepository {
	return &postgresAcademicUnitRepository{db: db}
}

func (r *postgresAcademicUnitRepository) Create(ctx context.Context, unit *entities.AcademicUnit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *postgresAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error) {
	var u entities.AcademicUnit
	query := r.db.WithContext(ctx)
	if includeDeleted {
		query = query.Unscoped()
	}
	if err := query.First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := r.db.WithContext(ctx)
	if includeDeleted {
		query = query.Unscoped()
	}
	var units []*entities.AcademicUnit
	err := query.Where("school_id = ?", schoolID).Order("created_at").Find(&units).Error
	return units, err
}

func (r *postgresAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	query := r.db.WithContext(ctx)
	if includeDeleted {
		query = query.Unscoped()
	}
	var units []*entities.AcademicUnit
	err := query.Where("school_id = ? AND type = ?", schoolID, unitType).Find(&units).Error
	return units, err
}

func (r *postgresAcademicUnitRepository) Update(ctx context.Context, unit *entities.AcademicUnit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

func (r *postgresAcademicUnitRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.AcademicUnit{}, "id = ?", id).Error
}

func (r *postgresAcademicUnitRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&entities.AcademicUnit{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": nil,
			"is_active":  true,
			"updated_at": time.Now(),
		}).Error
}

func (r *postgresAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error) {
	var units []*entities.AcademicUnit
	err := r.db.WithContext(ctx).Raw(`WITH RECURSIVE hierarchy AS (
		SELECT * FROM academic.academic_units WHERE id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT au.* FROM academic.academic_units au INNER JOIN hierarchy h ON au.id = h.parent_unit_id WHERE au.deleted_at IS NULL
	) SELECT * FROM hierarchy`, id).Scan(&units).Error
	if err != nil {
		return nil, err
	}
	// Reverse to get root-first order
	for i, j := 0, len(units)-1; i < j; i, j = i+1, j-1 {
		units[i], units[j] = units[j], units[i]
	}
	return units, nil
}

func (r *postgresAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.AcademicUnit{}).Where("school_id = ? AND code = ?", schoolID, code).Count(&count).Error
	return count > 0, err
}

// ==================== Subject ====================

type postgresSubjectRepository struct{ db *gorm.DB }

func NewPostgresSubjectRepository(db *gorm.DB) repository.SubjectRepository {
	return &postgresSubjectRepository{db: db}
}

func (r *postgresSubjectRepository) Create(ctx context.Context, s *entities.Subject) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *postgresSubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error) {
	var s entities.Subject
	if err := r.db.WithContext(ctx).Where("is_active = true").First(&s, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresSubjectRepository) Update(ctx context.Context, s *entities.Subject) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *postgresSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Subject{}).Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

func (r *postgresSubjectRepository) List(ctx context.Context) ([]*entities.Subject, error) {
	var subjects []*entities.Subject
	err := r.db.WithContext(ctx).Where("is_active = true").Order("name").Find(&subjects).Error
	return subjects, err
}

func (r *postgresSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Subject{}).Where("name = ? AND is_active = true", name).Count(&count).Error
	return count > 0, err
}

// ==================== Guardian ====================

type postgresGuardianRepository struct{ db *gorm.DB }

func NewPostgresGuardianRepository(db *gorm.DB) repository.GuardianRepository {
	return &postgresGuardianRepository{db: db}
}

func (r *postgresGuardianRepository) Create(ctx context.Context, g *entities.GuardianRelation) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *postgresGuardianRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	var g entities.GuardianRelation
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *postgresGuardianRepository) FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	var relations []*entities.GuardianRelation
	err := r.db.WithContext(ctx).Where("guardian_id = ? AND is_active = true", guardianID).Find(&relations).Error
	return relations, err
}

func (r *postgresGuardianRepository) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	var relations []*entities.GuardianRelation
	err := r.db.WithContext(ctx).Where("student_id = ? AND is_active = true", studentID).Find(&relations).Error
	return relations, err
}

func (r *postgresGuardianRepository) Update(ctx context.Context, g *entities.GuardianRelation) error {
	return r.db.WithContext(ctx).Save(g).Error
}

func (r *postgresGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.GuardianRelation{}).Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

func (r *postgresGuardianRepository) ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.GuardianRelation{}).
		Where("guardian_id = ? AND student_id = ? AND is_active = true", guardianID, studentID).Count(&count).Error
	return count > 0, err
}
