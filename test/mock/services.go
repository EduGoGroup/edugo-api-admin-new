package mock

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	authDTO "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/dto"
	"github.com/EduGoGroup/edugo-shared/auth"
)

// ---------------------------------------------------------------------------
// MockSchoolService
// ---------------------------------------------------------------------------

type MockSchoolService struct {
	CreateSchoolFn    func(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	GetSchoolFn       func(ctx context.Context, id string) (*dto.SchoolResponse, error)
	GetSchoolByCodeFn func(ctx context.Context, code string) (*dto.SchoolResponse, error)
	UpdateSchoolFn    func(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
	ListSchoolsFn     func(ctx context.Context) ([]dto.SchoolResponse, error)
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

func (m *MockSchoolService) ListSchools(ctx context.Context) ([]dto.SchoolResponse, error) {
	if m.ListSchoolsFn != nil {
		return m.ListSchoolsFn(ctx)
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
	CreateUnitFn       func(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	GetUnitFn          func(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	ListUnitsBySchoolFn func(ctx context.Context, schoolID string) ([]dto.AcademicUnitResponse, error)
	GetUnitTreeFn      func(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error)
	ListUnitsByTypeFn  func(ctx context.Context, schoolID, unitType string) ([]dto.AcademicUnitResponse, error)
	UpdateUnitFn       func(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	DeleteUnitFn       func(ctx context.Context, id string) error
	RestoreUnitFn      func(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	GetHierarchyPathFn func(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error)
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

func (m *MockAcademicUnitService) ListUnitsBySchool(ctx context.Context, schoolID string) ([]dto.AcademicUnitResponse, error) {
	if m.ListUnitsBySchoolFn != nil {
		return m.ListUnitsBySchoolFn(ctx, schoolID)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error) {
	if m.GetUnitTreeFn != nil {
		return m.GetUnitTreeFn(ctx, schoolID)
	}
	return nil, nil
}

func (m *MockAcademicUnitService) ListUnitsByType(ctx context.Context, schoolID, unitType string) ([]dto.AcademicUnitResponse, error) {
	if m.ListUnitsByTypeFn != nil {
		return m.ListUnitsByTypeFn(ctx, schoolID, unitType)
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
	ListMembershipsByUnitFn func(ctx context.Context, unitID string) ([]dto.MembershipResponse, error)
	ListMembershipsByRoleFn func(ctx context.Context, unitID, role string) ([]dto.MembershipResponse, error)
	ListMembershipsByUserFn func(ctx context.Context, userID string) ([]dto.MembershipResponse, error)
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

func (m *MockMembershipService) ListMembershipsByUnit(ctx context.Context, unitID string) ([]dto.MembershipResponse, error) {
	if m.ListMembershipsByUnitFn != nil {
		return m.ListMembershipsByUnitFn(ctx, unitID)
	}
	return nil, nil
}

func (m *MockMembershipService) ListMembershipsByRole(ctx context.Context, unitID, role string) ([]dto.MembershipResponse, error) {
	if m.ListMembershipsByRoleFn != nil {
		return m.ListMembershipsByRoleFn(ctx, unitID, role)
	}
	return nil, nil
}

func (m *MockMembershipService) ListMembershipsByUser(ctx context.Context, userID string) ([]dto.MembershipResponse, error) {
	if m.ListMembershipsByUserFn != nil {
		return m.ListMembershipsByUserFn(ctx, userID)
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
// MockRoleService
// ---------------------------------------------------------------------------

type MockRoleService struct {
	GetRolesFn           func(ctx context.Context, scope string) (*dto.RolesResponse, error)
	GetRoleFn            func(ctx context.Context, id string) (*dto.RoleDTO, error)
	GetRolePermissionsFn func(ctx context.Context, roleID string) (*dto.PermissionsResponse, error)
	GetUserRolesFn       func(ctx context.Context, userID string) (*dto.UserRolesResponse, error)
	GrantRoleToUserFn    func(ctx context.Context, userID string, req *dto.GrantRoleRequest, grantedBy string) (*dto.GrantRoleResponse, error)
	RevokeRoleFromUserFn func(ctx context.Context, userID, roleID string) error
}

func (m *MockRoleService) GetRoles(ctx context.Context, scope string) (*dto.RolesResponse, error) {
	if m.GetRolesFn != nil {
		return m.GetRolesFn(ctx, scope)
	}
	return nil, nil
}

func (m *MockRoleService) GetRole(ctx context.Context, id string) (*dto.RoleDTO, error) {
	if m.GetRoleFn != nil {
		return m.GetRoleFn(ctx, id)
	}
	return nil, nil
}

func (m *MockRoleService) GetRolePermissions(ctx context.Context, roleID string) (*dto.PermissionsResponse, error) {
	if m.GetRolePermissionsFn != nil {
		return m.GetRolePermissionsFn(ctx, roleID)
	}
	return nil, nil
}

func (m *MockRoleService) GetUserRoles(ctx context.Context, userID string) (*dto.UserRolesResponse, error) {
	if m.GetUserRolesFn != nil {
		return m.GetUserRolesFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockRoleService) GrantRoleToUser(ctx context.Context, userID string, req *dto.GrantRoleRequest, grantedBy string) (*dto.GrantRoleResponse, error) {
	if m.GrantRoleToUserFn != nil {
		return m.GrantRoleToUserFn(ctx, userID, req, grantedBy)
	}
	return nil, nil
}

func (m *MockRoleService) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	if m.RevokeRoleFromUserFn != nil {
		return m.RevokeRoleFromUserFn(ctx, userID, roleID)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockResourceService
// ---------------------------------------------------------------------------

type MockResourceService struct {
	ListResourcesFn  func(ctx context.Context) (*dto.ResourcesResponse, error)
	GetResourceFn    func(ctx context.Context, id string) (*dto.ResourceDTO, error)
	CreateResourceFn func(ctx context.Context, req dto.CreateResourceRequest) (*dto.ResourceDTO, error)
	UpdateResourceFn func(ctx context.Context, id string, req dto.UpdateResourceRequest) (*dto.ResourceDTO, error)
}

func (m *MockResourceService) ListResources(ctx context.Context) (*dto.ResourcesResponse, error) {
	if m.ListResourcesFn != nil {
		return m.ListResourcesFn(ctx)
	}
	return nil, nil
}

func (m *MockResourceService) GetResource(ctx context.Context, id string) (*dto.ResourceDTO, error) {
	if m.GetResourceFn != nil {
		return m.GetResourceFn(ctx, id)
	}
	return nil, nil
}

func (m *MockResourceService) CreateResource(ctx context.Context, req dto.CreateResourceRequest) (*dto.ResourceDTO, error) {
	if m.CreateResourceFn != nil {
		return m.CreateResourceFn(ctx, req)
	}
	return nil, nil
}

func (m *MockResourceService) UpdateResource(ctx context.Context, id string, req dto.UpdateResourceRequest) (*dto.ResourceDTO, error) {
	if m.UpdateResourceFn != nil {
		return m.UpdateResourceFn(ctx, id, req)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockMenuService
// ---------------------------------------------------------------------------

type MockMenuService struct {
	GetMenuForUserFn func(ctx context.Context, permissions []string) (*dto.MenuResponse, error)
	GetFullMenuFn    func(ctx context.Context) (*dto.MenuResponse, error)
}

func (m *MockMenuService) GetMenuForUser(ctx context.Context, permissions []string) (*dto.MenuResponse, error) {
	if m.GetMenuForUserFn != nil {
		return m.GetMenuForUserFn(ctx, permissions)
	}
	return nil, nil
}

func (m *MockMenuService) GetFullMenu(ctx context.Context) (*dto.MenuResponse, error) {
	if m.GetFullMenuFn != nil {
		return m.GetFullMenuFn(ctx)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockPermissionService
// ---------------------------------------------------------------------------

type MockPermissionService struct {
	ListPermissionsFn func(ctx context.Context) (*dto.PermissionsResponse, error)
	GetPermissionFn   func(ctx context.Context, id string) (*dto.PermissionDTO, error)
}

func (m *MockPermissionService) ListPermissions(ctx context.Context) (*dto.PermissionsResponse, error) {
	if m.ListPermissionsFn != nil {
		return m.ListPermissionsFn(ctx)
	}
	return nil, nil
}

func (m *MockPermissionService) GetPermission(ctx context.Context, id string) (*dto.PermissionDTO, error) {
	if m.GetPermissionFn != nil {
		return m.GetPermissionFn(ctx, id)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockSubjectService
// ---------------------------------------------------------------------------

type MockSubjectService struct {
	CreateSubjectFn func(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error)
	GetSubjectFn    func(ctx context.Context, id string) (*dto.SubjectResponse, error)
	ListSubjectsFn  func(ctx context.Context) ([]dto.SubjectResponse, error)
	UpdateSubjectFn func(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error)
	DeleteSubjectFn func(ctx context.Context, id string) error
}

func (m *MockSubjectService) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	if m.CreateSubjectFn != nil {
		return m.CreateSubjectFn(ctx, req)
	}
	return nil, nil
}

func (m *MockSubjectService) GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error) {
	if m.GetSubjectFn != nil {
		return m.GetSubjectFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSubjectService) ListSubjects(ctx context.Context) ([]dto.SubjectResponse, error) {
	if m.ListSubjectsFn != nil {
		return m.ListSubjectsFn(ctx)
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
// MockScreenConfigService
// ---------------------------------------------------------------------------

type MockScreenConfigService struct {
	CreateTemplateFn       func(ctx context.Context, req *service.CreateTemplateRequest) (*service.ScreenTemplateDTO, error)
	GetTemplateFn          func(ctx context.Context, id string) (*service.ScreenTemplateDTO, error)
	ListTemplatesFn        func(ctx context.Context, filter service.TemplateFilter) ([]*service.ScreenTemplateDTO, int, error)
	UpdateTemplateFn       func(ctx context.Context, id string, req *service.UpdateTemplateRequest) (*service.ScreenTemplateDTO, error)
	DeleteTemplateFn       func(ctx context.Context, id string) error
	CreateInstanceFn       func(ctx context.Context, req *service.CreateInstanceRequest) (*service.ScreenInstanceDTO, error)
	GetInstanceFn          func(ctx context.Context, id string) (*service.ScreenInstanceDTO, error)
	GetInstanceByKeyFn     func(ctx context.Context, key string) (*service.ScreenInstanceDTO, error)
	ListInstancesFn        func(ctx context.Context, filter service.InstanceFilter) ([]*service.ScreenInstanceDTO, int, error)
	UpdateInstanceFn       func(ctx context.Context, id string, req *service.UpdateInstanceRequest) (*service.ScreenInstanceDTO, error)
	DeleteInstanceFn       func(ctx context.Context, id string) error
	ResolveScreenByKeyFn   func(ctx context.Context, key string) (*service.CombinedScreenDTO, error)
	LinkScreenToResourceFn func(ctx context.Context, req *service.LinkScreenRequest) (*service.ResourceScreenDTO, error)
	GetScreensForResourceFn func(ctx context.Context, resourceID string) ([]*service.ResourceScreenDTO, error)
	UnlinkScreenFn         func(ctx context.Context, id string) error
}

func (m *MockScreenConfigService) CreateTemplate(ctx context.Context, req *service.CreateTemplateRequest) (*service.ScreenTemplateDTO, error) {
	if m.CreateTemplateFn != nil {
		return m.CreateTemplateFn(ctx, req)
	}
	return nil, nil
}

func (m *MockScreenConfigService) GetTemplate(ctx context.Context, id string) (*service.ScreenTemplateDTO, error) {
	if m.GetTemplateFn != nil {
		return m.GetTemplateFn(ctx, id)
	}
	return nil, nil
}

func (m *MockScreenConfigService) ListTemplates(ctx context.Context, filter service.TemplateFilter) ([]*service.ScreenTemplateDTO, int, error) {
	if m.ListTemplatesFn != nil {
		return m.ListTemplatesFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockScreenConfigService) UpdateTemplate(ctx context.Context, id string, req *service.UpdateTemplateRequest) (*service.ScreenTemplateDTO, error) {
	if m.UpdateTemplateFn != nil {
		return m.UpdateTemplateFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockScreenConfigService) DeleteTemplate(ctx context.Context, id string) error {
	if m.DeleteTemplateFn != nil {
		return m.DeleteTemplateFn(ctx, id)
	}
	return nil
}

func (m *MockScreenConfigService) CreateInstance(ctx context.Context, req *service.CreateInstanceRequest) (*service.ScreenInstanceDTO, error) {
	if m.CreateInstanceFn != nil {
		return m.CreateInstanceFn(ctx, req)
	}
	return nil, nil
}

func (m *MockScreenConfigService) GetInstance(ctx context.Context, id string) (*service.ScreenInstanceDTO, error) {
	if m.GetInstanceFn != nil {
		return m.GetInstanceFn(ctx, id)
	}
	return nil, nil
}

func (m *MockScreenConfigService) GetInstanceByKey(ctx context.Context, key string) (*service.ScreenInstanceDTO, error) {
	if m.GetInstanceByKeyFn != nil {
		return m.GetInstanceByKeyFn(ctx, key)
	}
	return nil, nil
}

func (m *MockScreenConfigService) ListInstances(ctx context.Context, filter service.InstanceFilter) ([]*service.ScreenInstanceDTO, int, error) {
	if m.ListInstancesFn != nil {
		return m.ListInstancesFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockScreenConfigService) UpdateInstance(ctx context.Context, id string, req *service.UpdateInstanceRequest) (*service.ScreenInstanceDTO, error) {
	if m.UpdateInstanceFn != nil {
		return m.UpdateInstanceFn(ctx, id, req)
	}
	return nil, nil
}

func (m *MockScreenConfigService) DeleteInstance(ctx context.Context, id string) error {
	if m.DeleteInstanceFn != nil {
		return m.DeleteInstanceFn(ctx, id)
	}
	return nil
}

func (m *MockScreenConfigService) ResolveScreenByKey(ctx context.Context, key string) (*service.CombinedScreenDTO, error) {
	if m.ResolveScreenByKeyFn != nil {
		return m.ResolveScreenByKeyFn(ctx, key)
	}
	return nil, nil
}

func (m *MockScreenConfigService) LinkScreenToResource(ctx context.Context, req *service.LinkScreenRequest) (*service.ResourceScreenDTO, error) {
	if m.LinkScreenToResourceFn != nil {
		return m.LinkScreenToResourceFn(ctx, req)
	}
	return nil, nil
}

func (m *MockScreenConfigService) GetScreensForResource(ctx context.Context, resourceID string) ([]*service.ResourceScreenDTO, error) {
	if m.GetScreensForResourceFn != nil {
		return m.GetScreensForResourceFn(ctx, resourceID)
	}
	return nil, nil
}

func (m *MockScreenConfigService) UnlinkScreen(ctx context.Context, id string) error {
	if m.UnlinkScreenFn != nil {
		return m.UnlinkScreenFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockAuthService
// ---------------------------------------------------------------------------

type MockAuthService struct {
	LoginFn                 func(ctx context.Context, email, password string) (*authDTO.LoginResponse, error)
	LogoutFn                func(ctx context.Context, accessToken string) error
	RefreshTokenFn          func(ctx context.Context, refreshToken string) (*authDTO.RefreshResponse, error)
	SwitchContextFn         func(ctx context.Context, userID, targetSchoolID string) (*authDTO.SwitchContextResponse, error)
	GetAvailableContextsFn  func(ctx context.Context, userID string, currentContext *auth.UserContext) (*authDTO.AvailableContextsResponse, error)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*authDTO.LoginResponse, error) {
	if m.LoginFn != nil {
		return m.LoginFn(ctx, email, password)
	}
	return nil, nil
}

func (m *MockAuthService) Logout(ctx context.Context, accessToken string) error {
	if m.LogoutFn != nil {
		return m.LogoutFn(ctx, accessToken)
	}
	return nil
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*authDTO.RefreshResponse, error) {
	if m.RefreshTokenFn != nil {
		return m.RefreshTokenFn(ctx, refreshToken)
	}
	return nil, nil
}

func (m *MockAuthService) SwitchContext(ctx context.Context, userID, targetSchoolID string) (*authDTO.SwitchContextResponse, error) {
	if m.SwitchContextFn != nil {
		return m.SwitchContextFn(ctx, userID, targetSchoolID)
	}
	return nil, nil
}

func (m *MockAuthService) GetAvailableContexts(ctx context.Context, userID string, currentContext *auth.UserContext) (*authDTO.AvailableContextsResponse, error) {
	if m.GetAvailableContextsFn != nil {
		return m.GetAvailableContextsFn(ctx, userID, currentContext)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MockUserService
// ---------------------------------------------------------------------------

type MockUserService struct {
	CreateUserFn func(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserFn    func(ctx context.Context, id string) (*dto.UserResponse, error)
	ListUsersFn  func(ctx context.Context, filters repository.ListFilters) ([]*dto.UserResponse, error)
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

func (m *MockUserService) ListUsers(ctx context.Context, filters repository.ListFilters) ([]*dto.UserResponse, error) {
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
