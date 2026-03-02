package mock

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

// ---------------------------------------------------------------------------
// MockSchoolService
// ---------------------------------------------------------------------------

type MockSchoolService struct {
	CreateSchoolFn    func(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	GetSchoolFn       func(ctx context.Context, id string) (*dto.SchoolResponse, error)
	GetSchoolByCodeFn func(ctx context.Context, code string) (*dto.SchoolResponse, error)
	UpdateSchoolFn    func(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
	ListSchoolsFn     func(ctx context.Context, filters sharedrepo.ListFilters) ([]dto.SchoolResponse, error)
	DeleteSchoolFn    func(ctx context.Context, id string) error
}

func (m *MockSchoolService) CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
	if m.CreateSchoolFn != nil {
		return m.CreateSchoolFn(ctx, req)
	}
	return nil, nil
}

func (m *MockSchoolService) GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error) {
	if m.GetSchoolFn != nil {
		return m.GetSchoolFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSchoolService) GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error) {
	if m.GetSchoolByCodeFn != nil {
		return m.GetSchoolByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *MockSchoolService) UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error) {
	if m.UpdateSchoolFn != nil {
		return m.UpdateSchoolFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockSchoolService) ListSchools(ctx context.Context, filters sharedrepo.ListFilters) ([]dto.SchoolResponse, error) {
	if m.ListSchoolsFn != nil {
		return m.ListSchoolsFn(ctx, filters)
	}
	return nil, nil
}

func (m *MockSchoolService) DeleteSchool(ctx context.Context, id string) error {
	if m.DeleteSchoolFn != nil {
		return m.DeleteSchoolFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockAcademicUnitService
// ---------------------------------------------------------------------------

type MockAcademicUnitService struct {
	CreateUnitFn        func(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	GetUnitFn           func(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	ListUnitsBySchoolFn func(ctx context.Context, schoolID string, filters sharedrepo.ListFilters) ([]dto.AcademicUnitResponse, error)
	GetUnitTreeFn       func(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error)
	ListUnitsByTypeFn   func(ctx context.Context, schoolID, unitType string, filters sharedrepo.ListFilters) ([]dto.AcademicUnitResponse, error)
	UpdateUnitFn        func(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	DeleteUnitFn        func(ctx context.Context, id string) error
	RestoreUnitFn       func(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	GetHierarchyPathFn  func(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error)
}

func (m *MockAcademicUnitService) CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	if m.CreateUnitFn != nil {
		return m.CreateUnitFn(ctx, schoolID, req)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	if m.GetUnitFn != nil {
		return m.GetUnitFn(ctx, id)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) ListUnitsBySchool(ctx context.Context, schoolID string, filters sharedrepo.ListFilters) ([]dto.AcademicUnitResponse, error) {
	if m.ListUnitsBySchoolFn != nil {
		return m.ListUnitsBySchoolFn(ctx, schoolID, filters)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error) {
	if m.GetUnitTreeFn != nil {
		return m.GetUnitTreeFn(ctx, schoolID)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) ListUnitsByType(ctx context.Context, schoolID, unitType string, filters sharedrepo.ListFilters) ([]dto.AcademicUnitResponse, error) {
	if m.ListUnitsByTypeFn != nil {
		return m.ListUnitsByTypeFn(ctx, schoolID, unitType, filters)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	if m.UpdateUnitFn != nil {
		return m.UpdateUnitFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) DeleteUnit(ctx context.Context, id string) error {
	if m.DeleteUnitFn != nil {
		return m.DeleteUnitFn(ctx, id)
	}
	return nil
}

func (m *MockAcademicUnitService) RestoreUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	if m.RestoreUnitFn != nil {
		return m.RestoreUnitFn(ctx, id)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error) {
	if m.GetHierarchyPathFn != nil {
		return m.GetHierarchyPathFn(ctx, id)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockMembershipService
// ---------------------------------------------------------------------------

type MockMembershipService struct {
	CreateMembershipFn      func(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error)
	GetMembershipFn         func(ctx context.Context, id string) (*dto.MembershipResponse, error)
	ListMembershipsByUnitFn func(ctx context.Context, unitID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error)
	ListMembershipsByRoleFn func(ctx context.Context, unitID, role string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error)
	ListMembershipsByUserFn func(ctx context.Context, userID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error)
	UpdateMembershipFn      func(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error)
	DeleteMembershipFn      func(ctx context.Context, id string) error
	ExpireMembershipFn      func(ctx context.Context, id string) (*dto.MembershipResponse, error)
}

func (m *MockMembershipService) CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error) {
	if m.CreateMembershipFn != nil {
		return m.CreateMembershipFn(ctx, req)
	}
	return nil, nil
}

func (m *MockMembershipService) GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	if m.GetMembershipFn != nil {
		return m.GetMembershipFn(ctx, id)
	}
	return nil, nil
}

func (m *MockMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error) {
	if m.ListMembershipsByUnitFn != nil {
		return m.ListMembershipsByUnitFn(ctx, unitID, filters)
	}
	return nil, nil
}

func (m *MockMembershipService) ListMembershipsByRole(ctx context.Context, unitID, role string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error) {
	if m.ListMembershipsByRoleFn != nil {
		return m.ListMembershipsByRoleFn(ctx, unitID, role, filters)
	}
	return nil, nil
}

func (m *MockMembershipService) ListMembershipsByUser(ctx context.Context, userID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error) {
	if m.ListMembershipsByUserFn != nil {
		return m.ListMembershipsByUserFn(ctx, userID, filters)
	}
	return nil, nil
}

func (m *MockMembershipService) UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error) {
	if m.UpdateMembershipFn != nil {
		return m.UpdateMembershipFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockMembershipService) DeleteMembership(ctx context.Context, id string) error {
	if m.DeleteMembershipFn != nil {
		return m.DeleteMembershipFn(ctx, id)
	}
	return nil
}

func (m *MockMembershipService) ExpireMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	if m.ExpireMembershipFn != nil {
		return m.ExpireMembershipFn(ctx, id)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockSubjectService
// ---------------------------------------------------------------------------

type MockSubjectService struct {
	CreateSubjectFn func(ctx context.Context, schoolID string, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error)
	GetSubjectFn    func(ctx context.Context, id string) (*dto.SubjectResponse, error)
	ListSubjectsFn  func(ctx context.Context, schoolID string, filters sharedrepo.ListFilters) ([]dto.SubjectResponse, error)
	UpdateSubjectFn func(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error)
	DeleteSubjectFn func(ctx context.Context, id string) error
}

func (m *MockSubjectService) CreateSubject(ctx context.Context, schoolID string, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	if m.CreateSubjectFn != nil {
		return m.CreateSubjectFn(ctx, schoolID, req)
	}
	return nil, nil
}

func (m *MockSubjectService) GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error) {
	if m.GetSubjectFn != nil {
		return m.GetSubjectFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSubjectService) ListSubjects(ctx context.Context, schoolID string, filters sharedrepo.ListFilters) ([]dto.SubjectResponse, error) {
	if m.ListSubjectsFn != nil {
		return m.ListSubjectsFn(ctx, schoolID, filters)
	}
	return nil, nil
}

func (m *MockSubjectService) UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	if m.UpdateSubjectFn != nil {
		return m.UpdateSubjectFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockSubjectService) DeleteSubject(ctx context.Context, id string) error {
	if m.DeleteSubjectFn != nil {
		return m.DeleteSubjectFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockGuardianService
// ---------------------------------------------------------------------------

type MockGuardianService struct {
	CreateRelationFn       func(ctx context.Context, req dto.CreateGuardianRelationRequest, createdBy string) (*dto.GuardianRelationResponse, error)
	GetRelationFn          func(ctx context.Context, id string) (*dto.GuardianRelationResponse, error)
	UpdateRelationFn       func(ctx context.Context, id string, req dto.UpdateGuardianRelationRequest) (*dto.GuardianRelationResponse, error)
	DeleteRelationFn       func(ctx context.Context, id string) error
	GetGuardianRelationsFn func(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error)
	GetStudentGuardiansFn  func(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error)
}

func (m *MockGuardianService) CreateRelation(ctx context.Context, req dto.CreateGuardianRelationRequest, createdBy string) (*dto.GuardianRelationResponse, error) {
	if m.CreateRelationFn != nil {
		return m.CreateRelationFn(ctx, req, createdBy)
	}
	return nil, nil
}

func (m *MockGuardianService) GetRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error) {
	if m.GetRelationFn != nil {
		return m.GetRelationFn(ctx, id)
	}
	return nil, nil
}

func (m *MockGuardianService) UpdateRelation(ctx context.Context, id string, req dto.UpdateGuardianRelationRequest) (*dto.GuardianRelationResponse, error) {
	if m.UpdateRelationFn != nil {
		return m.UpdateRelationFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockGuardianService) DeleteRelation(ctx context.Context, id string) error {
	if m.DeleteRelationFn != nil {
		return m.DeleteRelationFn(ctx, id)
	}
	return nil
}

func (m *MockGuardianService) GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error) {
	if m.GetGuardianRelationsFn != nil {
		return m.GetGuardianRelationsFn(ctx, guardianID)
	}
	return nil, nil
}

func (m *MockGuardianService) GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error) {
	if m.GetStudentGuardiansFn != nil {
		return m.GetStudentGuardiansFn(ctx, studentID)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockUserService
// ---------------------------------------------------------------------------

type MockUserService struct {
	CreateUserFn func(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserFn    func(ctx context.Context, id string) (*dto.UserResponse, error)
	ListUsersFn  func(ctx context.Context, filters sharedrepo.ListFilters) ([]*dto.UserResponse, error)
	UpdateUserFn func(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUserFn func(ctx context.Context, id string) error
}

func (m *MockUserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, req)
	}
	return nil, nil
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserService) ListUsers(ctx context.Context, filters sharedrepo.ListFilters) ([]*dto.UserResponse, error) {
	if m.ListUsersFn != nil {
		return m.ListUsersFn(ctx, filters)
	}
	return nil, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockStatsService
// ---------------------------------------------------------------------------

type MockStatsService struct {
	GetGlobalStatsFn func(ctx context.Context) (*dto.GlobalStatsResponse, error)
}

func (m *MockStatsService) GetGlobalStats(ctx context.Context) (*dto.GlobalStatsResponse, error) {
	if m.GetGlobalStatsFn != nil {
		return m.GetGlobalStatsFn(ctx)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockMaterialService
// ---------------------------------------------------------------------------

type MockMaterialService struct {
	DeleteMaterialFn func(ctx context.Context, id string) error
}

func (m *MockMaterialService) DeleteMaterial(ctx context.Context, id string) error {
	if m.DeleteMaterialFn != nil {
		return m.DeleteMaterialFn(ctx, id)
	}
	return nil
}
