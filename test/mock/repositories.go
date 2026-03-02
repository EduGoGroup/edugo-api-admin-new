package mock

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
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
	ListFn         func(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.School, error)
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

func (m *MockSchoolRepository) List(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.School, error) {
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
	FindBySchoolIDFn          func(ctx context.Context, schoolID uuid.UUID, includeDeleted bool, filters sharedrepo.ListFilters) ([]*entities.AcademicUnit, error)
	FindByTypeFn              func(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool, filters sharedrepo.ListFilters) ([]*entities.AcademicUnit, error)
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

func (m *MockAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, includeDeleted bool, filters sharedrepo.ListFilters) ([]*entities.AcademicUnit, error) {
	if m.FindBySchoolIDFn != nil {
		return m.FindBySchoolIDFn(ctx, schoolID, includeDeleted, filters)
	}
	return nil, nil
}

func (m *MockAcademicUnitRepository) FindByType(ctx context.Context, schoolID uuid.UUID, unitType string, includeDeleted bool, filters sharedrepo.ListFilters) ([]*entities.AcademicUnit, error) {
	if m.FindByTypeFn != nil {
		return m.FindByTypeFn(ctx, schoolID, unitType, includeDeleted, filters)
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
	CreateFn              func(ctx context.Context, membership *entities.Membership) error
	FindByIDFn            func(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	FindByUserFn          func(ctx context.Context, userID uuid.UUID, filters sharedrepo.ListFilters) ([]*entities.Membership, error)
	FindByUnitFn          func(ctx context.Context, unitID uuid.UUID, filters sharedrepo.ListFilters) ([]*entities.Membership, error)
	FindByUnitAndRoleFn   func(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool, filters sharedrepo.ListFilters) ([]*entities.Membership, error)
	FindByUserAndSchoolFn func(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error)
	UpdateFn              func(ctx context.Context, membership *entities.Membership) error
	DeleteFn              func(ctx context.Context, id uuid.UUID) error
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

func (m *MockMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID, filters sharedrepo.ListFilters) ([]*entities.Membership, error) {
	if m.FindByUserFn != nil {
		return m.FindByUserFn(ctx, userID, filters)
	}
	return nil, nil
}

func (m *MockMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID, filters sharedrepo.ListFilters) ([]*entities.Membership, error) {
	if m.FindByUnitFn != nil {
		return m.FindByUnitFn(ctx, unitID, filters)
	}
	return nil, nil
}

func (m *MockMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool, filters sharedrepo.ListFilters) ([]*entities.Membership, error) {
	if m.FindByUnitAndRoleFn != nil {
		return m.FindByUnitAndRoleFn(ctx, unitID, role, activeOnly, filters)
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
// MockSubjectRepository
// ---------------------------------------------------------------------------

type MockSubjectRepository struct {
	CreateFn                  func(ctx context.Context, subject *entities.Subject) error
	FindByIDFn                func(ctx context.Context, id uuid.UUID) (*entities.Subject, error)
	FindBySchoolIDFn          func(ctx context.Context, schoolID uuid.UUID, filters sharedrepo.ListFilters) ([]*entities.Subject, error)
	UpdateFn                  func(ctx context.Context, subject *entities.Subject) error
	DeleteFn                  func(ctx context.Context, id uuid.UUID) error
	ListFn                    func(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.Subject, error)
	ExistsByNameFn            func(ctx context.Context, name string) (bool, error)
	ExistsBySchoolIDAndNameFn func(ctx context.Context, schoolID uuid.UUID, name string) (bool, error)
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

func (m *MockSubjectRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, filters sharedrepo.ListFilters) ([]*entities.Subject, error) {
	if m.FindBySchoolIDFn != nil {
		return m.FindBySchoolIDFn(ctx, schoolID, filters)
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

func (m *MockSubjectRepository) List(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.Subject, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return nil, nil
}

func (m *MockSubjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.ExistsByNameFn != nil {
		return m.ExistsByNameFn(ctx, name)
	}
	return false, nil
}

func (m *MockSubjectRepository) ExistsBySchoolIDAndName(ctx context.Context, schoolID uuid.UUID, name string) (bool, error) {
	if m.ExistsBySchoolIDAndNameFn != nil {
		return m.ExistsBySchoolIDAndNameFn(ctx, schoolID, name)
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
	CreateFn        func(ctx context.Context, user *entities.User) error
	FindByIDFn      func(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmailFn   func(ctx context.Context, email string) (*entities.User, error)
	ExistsByEmailFn func(ctx context.Context, email string) (bool, error)
	UpdateFn        func(ctx context.Context, user *entities.User) error
	DeleteFn        func(ctx context.Context, id uuid.UUID) error
	ListFn          func(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.User, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
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

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailFn != nil {
		return m.ExistsByEmailFn(ctx, email)
	}
	return false, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, filters sharedrepo.ListFilters) ([]*entities.User, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return nil, nil
}
