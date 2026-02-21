package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleService_GetRoles(t *testing.T) {
	tests := []struct {
		name      string
		scope     string
		setupMock func(r *mock.MockRoleRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name:  "success - all roles",
			scope: "",
			setupMock: func(r *mock.MockRoleRepository) {
				r.FindAllFn = func(_ context.Context) ([]*entities.Role, error) {
					return []*entities.Role{
						{ID: uuid.New(), Name: "admin", DisplayName: "Admin", Scope: "system"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:  "success - filtered by scope",
			scope: "school",
			setupMock: func(r *mock.MockRoleRepository) {
				r.FindByScopeFn = func(_ context.Context, _ string) ([]*entities.Role, error) {
					return []*entities.Role{}, nil
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:  "error - database error",
			scope: "",
			setupMock: func(r *mock.MockRoleRepository) {
				r.FindAllFn = func(_ context.Context) ([]*entities.Role, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := &mock.MockRoleRepository{}
			if tt.setupMock != nil {
				tt.setupMock(roleRepo)
			}

			svc := service.NewRoleService(roleRepo, &mock.MockPermissionRepository{}, &mock.MockUserRoleRepository{}, mock.NewMockLogger())
			result, err := svc.GetRoles(context.Background(), tt.scope)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Roles, tt.wantCount)
			}
		})
	}
}

func TestRoleService_GetRole(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(r *mock.MockRoleRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(r *mock.MockRoleRepository) {
				r.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Role, error) {
					return &entities.Role{ID: validID, Name: "admin", DisplayName: "Admin", Scope: "system"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:        "error - invalid ID",
			id:          "bad",
			setupMock:   func(_ *mock.MockRoleRepository) {},
			wantErr:     true,
			errContains: "invalid role ID",
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(r *mock.MockRoleRepository) {
				r.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Role, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := &mock.MockRoleRepository{}
			if tt.setupMock != nil {
				tt.setupMock(roleRepo)
			}

			svc := service.NewRoleService(roleRepo, &mock.MockPermissionRepository{}, &mock.MockUserRoleRepository{}, mock.NewMockLogger())
			result, err := svc.GetRole(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestRoleService_GrantRoleToUser(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()

	tests := []struct {
		name        string
		userID      string
		request     *dto.GrantRoleRequest
		grantedBy   string
		setupMock   func(rr *mock.MockRoleRepository, ur *mock.MockUserRoleRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - grants role",
			userID: userID.String(),
			request: &dto.GrantRoleRequest{RoleID: roleID.String()},
			grantedBy: uuid.New().String(),
			setupMock: func(rr *mock.MockRoleRepository, ur *mock.MockUserRoleRepository) {
				rr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Role, error) {
					return &entities.Role{ID: roleID, Name: "teacher", DisplayName: "Teacher"}, nil
				}
				ur.UserHasRoleFn = func(_ context.Context, _, _ uuid.UUID, _, _ *uuid.UUID) (bool, error) {
					return false, nil
				}
				ur.GrantFn = func(_ context.Context, _ *entities.UserRole) error { return nil }
			},
			wantErr: false,
		},
		{
			name:        "error - invalid user ID",
			userID:      "bad",
			request:     &dto.GrantRoleRequest{RoleID: roleID.String()},
			setupMock:   func(_ *mock.MockRoleRepository, _ *mock.MockUserRoleRepository) {},
			wantErr:     true,
			errContains: "invalid user ID",
		},
		{
			name:   "error - role already granted",
			userID: userID.String(),
			request: &dto.GrantRoleRequest{RoleID: roleID.String()},
			setupMock: func(rr *mock.MockRoleRepository, ur *mock.MockUserRoleRepository) {
				rr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Role, error) {
					return &entities.Role{ID: roleID, Name: "teacher"}, nil
				}
				ur.UserHasRoleFn = func(_ context.Context, _, _ uuid.UUID, _, _ *uuid.UUID) (bool, error) {
					return true, nil
				}
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := &mock.MockRoleRepository{}
			userRoleRepo := &mock.MockUserRoleRepository{}
			if tt.setupMock != nil {
				tt.setupMock(roleRepo, userRoleRepo)
			}

			svc := service.NewRoleService(roleRepo, &mock.MockPermissionRepository{}, userRoleRepo, mock.NewMockLogger())
			result, err := svc.GrantRoleToUser(context.Background(), tt.userID, tt.request, tt.grantedBy)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestRoleService_RevokeRoleFromUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		roleID      string
		setupMock   func(ur *mock.MockUserRoleRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success",
			userID: uuid.New().String(),
			roleID: uuid.New().String(),
			setupMock: func(ur *mock.MockUserRoleRepository) {
				ur.RevokeByUserAndRoleFn = func(_ context.Context, _, _ uuid.UUID, _, _ *uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name:        "error - invalid user ID",
			userID:      "bad",
			roleID:      uuid.New().String(),
			setupMock:   func(_ *mock.MockUserRoleRepository) {},
			wantErr:     true,
			errContains: "invalid user ID",
		},
		{
			name:        "error - invalid role ID",
			userID:      uuid.New().String(),
			roleID:      "bad",
			setupMock:   func(_ *mock.MockUserRoleRepository) {},
			wantErr:     true,
			errContains: "invalid role ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRoleRepo := &mock.MockUserRoleRepository{}
			if tt.setupMock != nil {
				tt.setupMock(userRoleRepo)
			}

			svc := service.NewRoleService(&mock.MockRoleRepository{}, &mock.MockPermissionRepository{}, userRoleRepo, mock.NewMockLogger())
			err := svc.RevokeRoleFromUser(context.Background(), tt.userID, tt.roleID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRoleService_GetUserRoles(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()

	tests := []struct {
		name      string
		userID    string
		setupMock func(rr *mock.MockRoleRepository, ur *mock.MockUserRoleRepository)
		wantErr   bool
	}{
		{
			name:   "success",
			userID: userID.String(),
			setupMock: func(rr *mock.MockRoleRepository, ur *mock.MockUserRoleRepository) {
				ur.FindByUserFn = func(_ context.Context, _ uuid.UUID) ([]*entities.UserRole, error) {
					return []*entities.UserRole{
						{ID: uuid.New(), UserID: userID, RoleID: roleID, IsActive: true, GrantedAt: time.Now()},
					}, nil
				}
				rr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Role, error) {
					return &entities.Role{ID: roleID, Name: "admin"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid user ID",
			userID:    "bad",
			setupMock: func(_ *mock.MockRoleRepository, _ *mock.MockUserRoleRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := &mock.MockRoleRepository{}
			userRoleRepo := &mock.MockUserRoleRepository{}
			if tt.setupMock != nil {
				tt.setupMock(roleRepo, userRoleRepo)
			}

			svc := service.NewRoleService(roleRepo, &mock.MockPermissionRepository{}, userRoleRepo, mock.NewMockLogger())
			result, err := svc.GetUserRoles(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
