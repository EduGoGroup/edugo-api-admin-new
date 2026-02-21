package mock

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// MockSchoolRepository
// ---------------------------------------------------------------------------

type MockSchoolRepository struct {
	CreateFn       func(ctx context.Context, school *entities.School) error
	FindByIDFn     func(ctx context.Context, id uuid.UUID) (*entities.School, error)
	FindByCodeFn   func(ctx context.Context, code string) (*entities.School, error)
	UpdateFn       func(ctx context.Context, school *entities.School) error
	DeleteFn       func(ctx context.Context, id uuid.UUID) error
	ListFn         func(ctx context.Context, filters repository.ListFilters) ([]*entities.School, error)
	ExistsByCodeFn func(ctx context.Context, code string) (bool, error)
}

func (m *MockSchoolRepository) Create(ctx context.Context, school *entities.School) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, school)
	}
	return nil
}

func (m *MockSchoolRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSchoolRepository) FindByCode(ctx context.Context, code string) (*entities.School, error) {
	if m.FindByCodeFn != nil {
		return m.FindByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *MockSchoolRepository) Update(ctx context.Context, school *entities.School) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, school)
	}
	return nil
}

func (m *MockSchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockSchoolRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.School, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return nil, nil
}

func (m *MockSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	if m.ExistsByCodeFn != nil {
		return m.ExistsByCodeFn(ctx, code)
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// MockAcademicUnitRepository
// ---------------------------------------------------------------------------

type MockAcademicUnitRepository struct {
	CreateFn                  func(ctx context.Context, unit *entities.AcademicUnit) error
	FindByIDFn                func(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error)
	FindBySchoolIDFn          func(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error)
	FindByTypeFn              func(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error)
	UpdateFn                  func(ctx context.Context, unit *entities.AcademicUnit) error
	SoftDeleteFn              func(ctx context.Context, id uuid.UUID) error
	RestoreFn                 func(ctx context.Context, id uuid.UUID) error
	GetHierarchyPathFn        func(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error)
	ExistsBySchoolIDAndCodeFn func(ctx context.Context, schoolID uuid.UUID, code string) (bool, error)
}

func (m *MockAcademicUnitRepository) Create(ctx context.Context, unit *entities.AcademicUnit) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, unit)
	}
	return nil
}

func (m *MockAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID, includeDeleted bool) (*entities.AcademicUnit, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id, includeDeleted)
	}
	return nil, nil
}

func (m *MockAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	if m.FindBySchoolIDFn != nil {
		return m.FindBySchoolIDFn(ctx, schoolID, includeDeleted)
	}
	return nil, nil
}

func (m *MockAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool) ([]*entities.AcademicUnit, error) {
	if m.FindByTypeFn != nil {
		return m.FindByTypeFn(ctx, schoolID, unitType, includeDeleted)
	}
	return nil, nil
}

func (m *MockAcademicUnitRepository) Update(ctx context.Context, unit *entities.AcademicUnit) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, unit)
	}
	return nil
}

func (m *MockAcademicUnitRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	return nil
}

func (m *MockAcademicUnitRepository) Restore(ctx context.Context, id uuid.UUID) error {
	if m.RestoreFn != nil {
		return m.RestoreFn(ctx, id)
	}
	return nil
}

func (m *MockAcademicUnitRepository) GetHierarchyPath(ctx context.Context, id uuid.UUID) ([]*entities.AcademicUnit, error) {
	if m.GetHierarchyPathFn != nil {
		return m.GetHierarchyPathFn(ctx, id)
	}
	return nil, nil
}

func (m *MockAcademicUnitRepository) ExistsBySchoolIDAndCode(ctx context.Context, schoolID uuid.UUID, code string) (bool, error) {
	if m.ExistsBySchoolIDAndCodeFn != nil {
		return m.ExistsBySchoolIDAndCodeFn(ctx, schoolID, code)
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// MockMembershipRepository
// ---------------------------------------------------------------------------

type MockMembershipRepository struct {
	CreateFn             func(ctx context.Context, membership *entities.Membership) error
	FindByIDFn           func(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	FindByUserFn         func(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error)
	FindByUnitFn         func(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error)
	FindByUnitAndRoleFn  func(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error)
	FindByUserAndSchoolFn func(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error)
	UpdateFn             func(ctx context.Context, membership *entities.Membership) error
	DeleteFn             func(ctx context.Context, id uuid.UUID) error
}

func (m *MockMembershipRepository) Create(ctx context.Context, membership *entities.Membership) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, membership)
	}
	return nil
}

func (m *MockMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error) {
	if m.FindByUserFn != nil {
		return m.FindByUserFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error) {
	if m.FindByUnitFn != nil {
		return m.FindByUnitFn(ctx, unitID)
	}
	return nil, nil
}

func (m *MockMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error) {
	if m.FindByUnitAndRoleFn != nil {
		return m.FindByUnitAndRoleFn(ctx, unitID, role, activeOnly)
	}
	return nil, nil
}

func (m *MockMembershipRepository) FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error) {
	if m.FindByUserAndSchoolFn != nil {
		return m.FindByUserAndSchoolFn(ctx, userID, schoolID)
	}
	return nil, nil
}

func (m *MockMembershipRepository) Update(ctx context.Context, membership *entities.Membership) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, membership)
	}
	return nil
}

func (m *MockMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockRoleRepository
// ---------------------------------------------------------------------------

type MockRoleRepository struct {
	FindByIDFn    func(ctx context.Context, id uuid.UUID) (*entities.Role, error)
	FindAllFn     func(ctx context.Context) ([]*entities.Role, error)
	FindByScopeFn func(ctx context.Context, scope string) ([]*entities.Role, error)
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockRoleRepository) FindAll(ctx context.Context) ([]*entities.Role, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, nil
}

func (m *MockRoleRepository) FindByScope(ctx context.Context, scope string) ([]*entities.Role, error) {
	if m.FindByScopeFn != nil {
		return m.FindByScopeFn(ctx, scope)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockPermissionRepository
// ---------------------------------------------------------------------------

type MockPermissionRepository struct {
	FindByIDFn   func(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
	FindAllFn    func(ctx context.Context) ([]*entities.Permission, error)
	FindByRoleFn func(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error)
}

func (m *MockPermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockPermissionRepository) FindAll(ctx context.Context) ([]*entities.Permission, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, nil
}

func (m *MockPermissionRepository) FindByRole(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error) {
	if m.FindByRoleFn != nil {
		return m.FindByRoleFn(ctx, roleID)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockUserRoleRepository
// ---------------------------------------------------------------------------

type MockUserRoleRepository struct {
	FindByUserFn           func(ctx context.Context, userID uuid.UUID) ([]*entities.UserRole, error)
	FindByUserInContextFn  func(ctx context.Context, userID uuid.UUID, schoolID *uuid.UUID, unitID *uuid.UUID) ([]*entities.UserRole, error)
	GrantFn                func(ctx context.Context, userRole *entities.UserRole) error
	RevokeFn               func(ctx context.Context, id uuid.UUID) error
	RevokeByUserAndRoleFn  func(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) error
	UserHasRoleFn          func(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) (bool, error)
	GetUserPermissionsFn   func(ctx context.Context, userID uuid.UUID, schoolID, unitID *uuid.UUID) ([]string, error)
}

func (m *MockUserRoleRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.UserRole, error) {
	if m.FindByUserFn != nil {
		return m.FindByUserFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockUserRoleRepository) FindByUserInContext(ctx context.Context, userID uuid.UUID, schoolID *uuid.UUID, unitID *uuid.UUID) ([]*entities.UserRole, error) {
	if m.FindByUserInContextFn != nil {
		return m.FindByUserInContextFn(ctx, userID, schoolID, unitID)
	}
	return nil, nil
}

func (m *MockUserRoleRepository) Grant(ctx context.Context, userRole *entities.UserRole) error {
	if m.GrantFn != nil {
		return m.GrantFn(ctx, userRole)
	}
	return nil
}

func (m *MockUserRoleRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	if m.RevokeFn != nil {
		return m.RevokeFn(ctx, id)
	}
	return nil
}

func (m *MockUserRoleRepository) RevokeByUserAndRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) error {
	if m.RevokeByUserAndRoleFn != nil {
		return m.RevokeByUserAndRoleFn(ctx, userID, roleID, schoolID, unitID)
	}
	return nil
}

func (m *MockUserRoleRepository) UserHasRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) (bool, error) {
	if m.UserHasRoleFn != nil {
		return m.UserHasRoleFn(ctx, userID, roleID, schoolID, unitID)
	}
	return false, nil
}

func (m *MockUserRoleRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID, schoolID, unitID *uuid.UUID) ([]string, error) {
	if m.GetUserPermissionsFn != nil {
		return m.GetUserPermissionsFn(ctx, userID, schoolID, unitID)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockResourceRepository
// ---------------------------------------------------------------------------

type MockResourceRepository struct {
	FindAllFn        func(ctx context.Context) ([]*entities.Resource, error)
	FindByIDFn       func(ctx context.Context, id uuid.UUID) (*entities.Resource, error)
	FindMenuVisibleFn func(ctx context.Context) ([]*entities.Resource, error)
	CreateFn         func(ctx context.Context, resource *entities.Resource) error
	UpdateFn         func(ctx context.Context, resource *entities.Resource) error
}

func (m *MockResourceRepository) FindAll(ctx context.Context) ([]*entities.Resource, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, nil
}

func (m *MockResourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Resource, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockResourceRepository) FindMenuVisible(ctx context.Context) ([]*entities.Resource, error) {
	if m.FindMenuVisibleFn != nil {
		return m.FindMenuVisibleFn(ctx)
	}
	return nil, nil
}

func (m *MockResourceRepository) Create(ctx context.Context, resource *entities.Resource) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, resource)
	}
	return nil
}

func (m *MockResourceRepository) Update(ctx context.Context, resource *entities.Resource) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, resource)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockSubjectRepository
// ---------------------------------------------------------------------------

type MockSubjectRepository struct {
	CreateFn       func(ctx context.Context, subject *entities.Subject) error
	FindByIDFn     func(ctx context.Context, id uuid.UUID) (*entities.Subject, error)
	UpdateFn       func(ctx context.Context, subject *entities.Subject) error
	DeleteFn       func(ctx context.Context, id uuid.UUID) error
	ListFn         func(ctx context.Context) ([]*entities.Subject, error)
	ExistsByNameFn func(ctx context.Context, name string) (bool, error)
}

func (m *MockSubjectRepository) Create(ctx context.Context, subject *entities.Subject) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, subject)
	}
	return nil
}

func (m *MockSubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Subject, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSubjectRepository) Update(ctx context.Context, subject *entities.Subject) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, subject)
	}
	return nil
}

func (m *MockSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockSubjectRepository) List(ctx context.Context) ([]*entities.Subject, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return nil, nil
}

func (m *MockSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.ExistsByNameFn != nil {
		return m.ExistsByNameFn(ctx, name)
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// MockGuardianRepository
// ---------------------------------------------------------------------------

type MockGuardianRepository struct {
	CreateFn               func(ctx context.Context, relation *entities.GuardianRelation) error
	FindByIDFn             func(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error)
	FindByGuardianFn       func(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error)
	FindByStudentFn        func(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error)
	UpdateFn               func(ctx context.Context, relation *entities.GuardianRelation) error
	DeleteFn               func(ctx context.Context, id uuid.UUID) error
	ExistsActiveRelationFn func(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error)
}

func (m *MockGuardianRepository) Create(ctx context.Context, relation *entities.GuardianRelation) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, relation)
	}
	return nil
}

func (m *MockGuardianRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.GuardianRelation, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockGuardianRepository) FindByGuardian(ctx context.Context, guardianID uuid.UUID) ([]*entities.GuardianRelation, error) {
	if m.FindByGuardianFn != nil {
		return m.FindByGuardianFn(ctx, guardianID)
	}
	return nil, nil
}

func (m *MockGuardianRepository) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]*entities.GuardianRelation, error) {
	if m.FindByStudentFn != nil {
		return m.FindByStudentFn(ctx, studentID)
	}
	return nil, nil
}

func (m *MockGuardianRepository) Update(ctx context.Context, relation *entities.GuardianRelation) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, relation)
	}
	return nil
}

func (m *MockGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockGuardianRepository) ExistsActiveRelation(ctx context.Context, guardianID, studentID uuid.UUID) (bool, error) {
	if m.ExistsActiveRelationFn != nil {
		return m.ExistsActiveRelationFn(ctx, guardianID, studentID)
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// MockUserRepository
// ---------------------------------------------------------------------------

type MockUserRepository struct {
	FindByIDFn    func(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmailFn func(ctx context.Context, email string) (*entities.User, error)
	UpdateFn      func(ctx context.Context, user *entities.User) error
	ListFn        func(ctx context.Context, filters repository.ListFilters) ([]*entities.User, error)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entities.User, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockScreenTemplateRepository
// ---------------------------------------------------------------------------

type MockScreenTemplateRepository struct {
	CreateFn  func(ctx context.Context, template *entities.ScreenTemplate) error
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*entities.ScreenTemplate, error)
	ListFn    func(ctx context.Context, filter repository.ScreenTemplateFilter) ([]*entities.ScreenTemplate, int, error)
	UpdateFn  func(ctx context.Context, template *entities.ScreenTemplate) error
	DeleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *MockScreenTemplateRepository) Create(ctx context.Context, template *entities.ScreenTemplate) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, template)
	}
	return nil
}

func (m *MockScreenTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScreenTemplate, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockScreenTemplateRepository) List(ctx context.Context, filter repository.ScreenTemplateFilter) ([]*entities.ScreenTemplate, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockScreenTemplateRepository) Update(ctx context.Context, template *entities.ScreenTemplate) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, template)
	}
	return nil
}

func (m *MockScreenTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockScreenInstanceRepository
// ---------------------------------------------------------------------------

type MockScreenInstanceRepository struct {
	CreateFn         func(ctx context.Context, instance *entities.ScreenInstance) error
	GetByIDFn        func(ctx context.Context, id uuid.UUID) (*entities.ScreenInstance, error)
	GetByScreenKeyFn func(ctx context.Context, key string) (*entities.ScreenInstance, error)
	ListFn           func(ctx context.Context, filter repository.ScreenInstanceFilter) ([]*entities.ScreenInstance, int, error)
	UpdateFn         func(ctx context.Context, instance *entities.ScreenInstance) error
	DeleteFn         func(ctx context.Context, id uuid.UUID) error
}

func (m *MockScreenInstanceRepository) Create(ctx context.Context, instance *entities.ScreenInstance) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, instance)
	}
	return nil
}

func (m *MockScreenInstanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScreenInstance, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockScreenInstanceRepository) GetByScreenKey(ctx context.Context, key string) (*entities.ScreenInstance, error) {
	if m.GetByScreenKeyFn != nil {
		return m.GetByScreenKeyFn(ctx, key)
	}
	return nil, nil
}

func (m *MockScreenInstanceRepository) List(ctx context.Context, filter repository.ScreenInstanceFilter) ([]*entities.ScreenInstance, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockScreenInstanceRepository) Update(ctx context.Context, instance *entities.ScreenInstance) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, instance)
	}
	return nil
}

func (m *MockScreenInstanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockResourceScreenRepository
// ---------------------------------------------------------------------------

type MockResourceScreenRepository struct {
	CreateFn           func(ctx context.Context, rs *entities.ResourceScreen) error
	GetByResourceIDFn  func(ctx context.Context, resourceID uuid.UUID) ([]*entities.ResourceScreen, error)
	GetByResourceKeyFn func(ctx context.Context, key string) ([]*entities.ResourceScreen, error)
	DeleteFn           func(ctx context.Context, id uuid.UUID) error
}

func (m *MockResourceScreenRepository) Create(ctx context.Context, rs *entities.ResourceScreen) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, rs)
	}
	return nil
}

func (m *MockResourceScreenRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*entities.ResourceScreen, error) {
	if m.GetByResourceIDFn != nil {
		return m.GetByResourceIDFn(ctx, resourceID)
	}
	return nil, nil
}

func (m *MockResourceScreenRepository) GetByResourceKey(ctx context.Context, key string) ([]*entities.ResourceScreen, error) {
	if m.GetByResourceKeyFn != nil {
		return m.GetByResourceKeyFn(ctx, key)
	}
	return nil, nil
}

func (m *MockResourceScreenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
