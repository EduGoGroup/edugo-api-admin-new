package repository

import (
	"context"
	"errors"
	"fmt"
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

// ==================== Membership ====================

type postgresMembershipRepository struct{ db *gorm.DB }

func NewPostgresMembershipRepository(db *gorm.DB) repository.MembershipRepository {
	return &postgresMembershipRepository{db: db}
}

func (r *postgresMembershipRepository) Create(ctx context.Context, m *entities.Membership) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *postgresMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	var m entities.Membership
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *postgresMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error) {
	var memberships []*entities.Membership
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID).Order("created_at DESC").Find(&memberships).Error
	return memberships, err
}

func (r *postgresMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error) {
	var memberships []*entities.Membership
	err := r.db.WithContext(ctx).Where("academic_unit_id = ? AND is_active = true", unitID).Order("created_at DESC").Find(&memberships).Error
	return memberships, err
}

func (r *postgresMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error) {
	query := r.db.WithContext(ctx).Where("academic_unit_id = ? AND role = ?", unitID, role)
	if activeOnly {
		query = query.Where("is_active = true")
	}
	var memberships []*entities.Membership
	err := query.Find(&memberships).Error
	return memberships, err
}

func (r *postgresMembershipRepository) FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error) {
	var m entities.Membership
	if err := r.db.WithContext(ctx).Where("user_id = ? AND school_id = ? AND is_active = true", userID, schoolID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *postgresMembershipRepository) Update(ctx context.Context, m *entities.Membership) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *postgresMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Membership{}, "id = ?", id).Error
}

// ==================== Role ====================

type postgresRoleRepository struct{ db *gorm.DB }

func NewPostgresRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &postgresRoleRepository{db: db}
}

func (r *postgresRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.WithContext(ctx).Where("is_active = true").First(&role, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *postgresRoleRepository) FindAll(ctx context.Context) ([]*entities.Role, error) {
	var roles []*entities.Role
	err := r.db.WithContext(ctx).Where("is_active = true").Order("name").Find(&roles).Error
	return roles, err
}

func (r *postgresRoleRepository) FindByScope(ctx context.Context, scope string) ([]*entities.Role, error) {
	var roles []*entities.Role
	err := r.db.WithContext(ctx).Where("scope = ? AND is_active = true", scope).Order("name").Find(&roles).Error
	return roles, err
}

// ==================== Permission ====================

type postgresPermissionRepository struct{ db *gorm.DB }

func NewPostgresPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &postgresPermissionRepository{db: db}
}

func (r *postgresPermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	var p entities.Permission
	if err := r.db.WithContext(ctx).Where("is_active = true").First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresPermissionRepository) FindAll(ctx context.Context) ([]*entities.Permission, error) {
	var perms []*entities.Permission
	err := r.db.WithContext(ctx).Where("is_active = true").Order("name").Find(&perms).Error
	return perms, err
}

func (r *postgresPermissionRepository) FindByRole(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error) {
	var perms []*entities.Permission
	err := r.db.WithContext(ctx).
		Joins("INNER JOIN iam.role_permissions rp ON iam.permissions.id = rp.permission_id").
		Where("rp.role_id = ? AND iam.permissions.is_active = true", roleID).
		Order("iam.permissions.name").
		Find(&perms).Error
	return perms, err
}

// ==================== UserRole ====================

type postgresUserRoleRepository struct{ db *gorm.DB }

func NewPostgresUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return &postgresUserRoleRepository{db: db}
}

func (r *postgresUserRoleRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.UserRole, error) {
	var userRoles []*entities.UserRole
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID).Find(&userRoles).Error
	return userRoles, err
}

func (r *postgresUserRoleRepository) FindByUserInContext(ctx context.Context, userID uuid.UUID, schoolID *uuid.UUID, unitID *uuid.UUID) ([]*entities.UserRole, error) {
	query := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID)
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if unitID != nil {
		query = query.Where("academic_unit_id = ?", *unitID)
	}
	query = query.Order("created_at")
	var userRoles []*entities.UserRole
	err := query.Find(&userRoles).Error
	return userRoles, err
}

func (r *postgresUserRoleRepository) Grant(ctx context.Context, ur *entities.UserRole) error {
	return r.db.WithContext(ctx).Create(ur).Error
}

func (r *postgresUserRoleRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.UserRole{}).Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

func (r *postgresUserRoleRepository) RevokeByUserAndRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.UserRole{}).
		Where("user_id = ? AND role_id = ? AND is_active = true", userID, roleID).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

func (r *postgresUserRoleRepository) UserHasRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&entities.UserRole{}).Where("user_id = ? AND role_id = ? AND is_active = true", userID, roleID)
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if unitID != nil {
		query = query.Where("academic_unit_id = ?", *unitID)
	}
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *postgresUserRoleRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID, schoolID, unitID *uuid.UUID) ([]string, error) {
	query := `SELECT DISTINCT p.name FROM iam.permissions p
		INNER JOIN iam.role_permissions rp ON p.id = rp.permission_id
		INNER JOIN iam.user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND ur.is_active = true AND p.is_active = true`
	args := []interface{}{userID}
	if schoolID != nil {
		query += ` AND ur.school_id = ?`
		args = append(args, *schoolID)
	}
	if unitID != nil {
		query += ` AND ur.academic_unit_id = ?`
		args = append(args, *unitID)
	}
	var perms []string
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&perms).Error
	return perms, err
}

// ==================== Resource ====================

type postgresResourceRepository struct{ db *gorm.DB }

func NewPostgresResourceRepository(db *gorm.DB) repository.ResourceRepository {
	return &postgresResourceRepository{db: db}
}

func (r *postgresResourceRepository) FindAll(ctx context.Context) ([]*entities.Resource, error) {
	var resources []*entities.Resource
	err := r.db.WithContext(ctx).Where("is_active = true").Order("sort_order").Find(&resources).Error
	return resources, err
}

func (r *postgresResourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Resource, error) {
	var res entities.Resource
	if err := r.db.WithContext(ctx).Where("is_active = true").First(&res, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *postgresResourceRepository) FindMenuVisible(ctx context.Context) ([]*entities.Resource, error) {
	var resources []*entities.Resource
	err := r.db.WithContext(ctx).Where("is_menu_visible = true AND is_active = true").Order("sort_order").Find(&resources).Error
	return resources, err
}

func (r *postgresResourceRepository) Create(ctx context.Context, res *entities.Resource) error {
	return r.db.WithContext(ctx).Create(res).Error
}

func (r *postgresResourceRepository) Update(ctx context.Context, res *entities.Resource) error {
	return r.db.WithContext(ctx).Save(res).Error
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

// ==================== ScreenTemplate ====================

type postgresScreenTemplateRepository struct{ db *gorm.DB }

func NewPostgresScreenTemplateRepository(db *gorm.DB) repository.ScreenTemplateRepository {
	return &postgresScreenTemplateRepository{db: db}
}

func (r *postgresScreenTemplateRepository) Create(ctx context.Context, t *entities.ScreenTemplate) error {
	return r.db.WithContext(ctx).Table("ui_config.screen_templates").Create(t).Error
}

func (r *postgresScreenTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScreenTemplate, error) {
	var t entities.ScreenTemplate
	if err := r.db.WithContext(ctx).Table("ui_config.screen_templates").First(&t, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("screen template not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresScreenTemplateRepository) List(ctx context.Context, filter repository.ScreenTemplateFilter) ([]*entities.ScreenTemplate, int, error) {
	baseQuery := r.db.WithContext(ctx).Table("ui_config.screen_templates").Where("is_active = true")
	if filter.Pattern != "" {
		baseQuery = baseQuery.Where("pattern = ?", filter.Pattern)
	}

	var total int64
	baseQuery.Count(&total)

	query := r.db.WithContext(ctx).Table("ui_config.screen_templates").Where("is_active = true")
	if filter.Pattern != "" {
		query = query.Where("pattern = ?", filter.Pattern)
	}
	query = query.Order("created_at DESC")
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var templates []*entities.ScreenTemplate
	err := query.Find(&templates).Error
	return templates, int(total), err
}

func (r *postgresScreenTemplateRepository) Update(ctx context.Context, t *entities.ScreenTemplate) error {
	return r.db.WithContext(ctx).Table("ui_config.screen_templates").Save(t).Error
}

func (r *postgresScreenTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Table("ui_config.screen_templates").Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

// ==================== ScreenInstance ====================

type postgresScreenInstanceRepository struct{ db *gorm.DB }

func NewPostgresScreenInstanceRepository(db *gorm.DB) repository.ScreenInstanceRepository {
	return &postgresScreenInstanceRepository{db: db}
}

func (r *postgresScreenInstanceRepository) Create(ctx context.Context, i *entities.ScreenInstance) error {
	return r.db.WithContext(ctx).Table("ui_config.screen_instances").Create(i).Error
}

func (r *postgresScreenInstanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScreenInstance, error) {
	var i entities.ScreenInstance
	if err := r.db.WithContext(ctx).Table("ui_config.screen_instances").First(&i, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("screen instance not found")
		}
		return nil, err
	}
	return &i, nil
}

func (r *postgresScreenInstanceRepository) GetByScreenKey(ctx context.Context, key string) (*entities.ScreenInstance, error) {
	var i entities.ScreenInstance
	if err := r.db.WithContext(ctx).Table("ui_config.screen_instances").Where("screen_key = ? AND is_active = true", key).First(&i).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("screen instance not found for key: %s", key)
		}
		return nil, err
	}
	return &i, nil
}

func (r *postgresScreenInstanceRepository) List(ctx context.Context, filter repository.ScreenInstanceFilter) ([]*entities.ScreenInstance, int, error) {
	baseQuery := r.db.WithContext(ctx).Table("ui_config.screen_instances").Where("is_active = true")
	if filter.TemplateID != nil {
		baseQuery = baseQuery.Where("template_id = ?", *filter.TemplateID)
	}

	var total int64
	baseQuery.Count(&total)

	query := r.db.WithContext(ctx).Table("ui_config.screen_instances").Where("is_active = true")
	if filter.TemplateID != nil {
		query = query.Where("template_id = ?", *filter.TemplateID)
	}
	query = query.Order("created_at DESC")
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var instances []*entities.ScreenInstance
	err := query.Find(&instances).Error
	return instances, int(total), err
}

func (r *postgresScreenInstanceRepository) Update(ctx context.Context, i *entities.ScreenInstance) error {
	return r.db.WithContext(ctx).Table("ui_config.screen_instances").Save(i).Error
}

func (r *postgresScreenInstanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Table("ui_config.screen_instances").Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error
}

// ==================== ResourceScreen ====================

type postgresResourceScreenRepository struct{ db *gorm.DB }

func NewPostgresResourceScreenRepository(db *gorm.DB) repository.ResourceScreenRepository {
	return &postgresResourceScreenRepository{db: db}
}

func (r *postgresResourceScreenRepository) Create(ctx context.Context, rs *entities.ResourceScreen) error {
	return r.db.WithContext(ctx).Table("ui_config.resource_screens").Create(rs).Error
}

func (r *postgresResourceScreenRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*entities.ResourceScreen, error) {
	var result []*entities.ResourceScreen
	err := r.db.WithContext(ctx).Table("ui_config.resource_screens").
		Where("resource_id = ? AND is_active = true", resourceID).Order("sort_order").Find(&result).Error
	return result, err
}

func (r *postgresResourceScreenRepository) GetByResourceKey(ctx context.Context, key string) ([]*entities.ResourceScreen, error) {
	var result []*entities.ResourceScreen
	err := r.db.WithContext(ctx).Table("ui_config.resource_screens").
		Where("resource_key = ? AND is_active = true", key).Order("sort_order").Find(&result).Error
	return result, err
}

func (r *postgresResourceScreenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Table("ui_config.resource_screens").Where("id = ?", id).Delete(&entities.ResourceScreen{}).Error
}
