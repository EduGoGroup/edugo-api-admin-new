package service_test

import (
	"context"
	"testing"
	"time"

	authService "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/service"
	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTestDeps() (*mock.MockUserRepository, *mock.MockUserRoleRepository, *mock.MockRoleRepository, *mock.MockMembershipRepository, *mock.MockSchoolRepository, *authService.TokenService) {
	jwtManager := auth.NewJWTManager("test-secret-key-for-auth-testing", "edugo-test")
	tokenSvc := authService.NewTokenService(jwtManager, 15*time.Minute, 7*24*time.Hour)
	return &mock.MockUserRepository{}, &mock.MockUserRoleRepository{}, &mock.MockRoleRepository{}, &mock.MockMembershipRepository{}, &mock.MockSchoolRepository{}, tokenSvc
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := auth.HashPassword("correctpassword")
	schoolID := uuid.New()
	roleID := uuid.New()

	tests := []struct {
		name        string
		email       string
		password    string
		setupMock   func(ur *mock.MockUserRepository, urr *mock.MockUserRoleRepository, rr *mock.MockRoleRepository, mr *mock.MockMembershipRepository, sr *mock.MockSchoolRepository)
		wantErr     bool
		errTarget   error
	}{
		{
			name:     "success - valid login",
			email:    "admin@test.com",
			password: "correctpassword",
			setupMock: func(ur *mock.MockUserRepository, urr *mock.MockUserRoleRepository, rr *mock.MockRoleRepository, mr *mock.MockMembershipRepository, sr *mock.MockSchoolRepository) {
				ur.FindByEmailFn = func(_ context.Context, _ string) (*entities.User, error) {
					return &entities.User{
						ID: uuid.New(), Email: "admin@test.com", PasswordHash: hashedPassword,
						FirstName: "Admin", LastName: "User", IsActive: true,
					}, nil
				}
				urr.FindByUserInContextFn = func(_ context.Context, _ uuid.UUID, _ *uuid.UUID, _ *uuid.UUID) ([]*entities.UserRole, error) {
					return []*entities.UserRole{
						{ID: uuid.New(), UserID: uuid.New(), RoleID: roleID, IsActive: true, GrantedAt: time.Now()},
					}, nil
				}
				rr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Role, error) {
					return &entities.Role{ID: roleID, Name: "super_admin", DisplayName: "Super Admin"}, nil
				}
				mr.FindByUserFn = func(_ context.Context, _ uuid.UUID) ([]*entities.Membership, error) {
					return []*entities.Membership{
						{ID: uuid.New(), SchoolID: schoolID, IsActive: true},
					}, nil
				}
				sr.FindByIDFn = func(_ context.Context, id uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: id, Name: "Test School"}, nil
				}
				urr.GetUserPermissionsFn = func(_ context.Context, _ uuid.UUID, _, _ *uuid.UUID) ([]string, error) {
					return []string{"schools:read", "schools:create"}, nil
				}
				ur.UpdateFn = func(_ context.Context, _ *entities.User) error { return nil }
			},
			wantErr: false,
		},
		{
			name:     "error - user not found",
			email:    "nonexistent@test.com",
			password: "any",
			setupMock: func(ur *mock.MockUserRepository, _ *mock.MockUserRoleRepository, _ *mock.MockRoleRepository, _ *mock.MockMembershipRepository, _ *mock.MockSchoolRepository) {
				ur.FindByEmailFn = func(_ context.Context, _ string) (*entities.User, error) { return nil, nil }
			},
			wantErr:   true,
			errTarget: authService.ErrInvalidCredentials,
		},
		{
			name:     "error - user inactive",
			email:    "inactive@test.com",
			password: "any",
			setupMock: func(ur *mock.MockUserRepository, _ *mock.MockUserRoleRepository, _ *mock.MockRoleRepository, _ *mock.MockMembershipRepository, _ *mock.MockSchoolRepository) {
				ur.FindByEmailFn = func(_ context.Context, _ string) (*entities.User, error) {
					return &entities.User{
						ID: uuid.New(), Email: "inactive@test.com", PasswordHash: hashedPassword,
						IsActive: false,
					}, nil
				}
			},
			wantErr:   true,
			errTarget: authService.ErrUserInactive,
		},
		{
			name:     "error - wrong password",
			email:    "admin@test.com",
			password: "wrongpassword",
			setupMock: func(ur *mock.MockUserRepository, _ *mock.MockUserRoleRepository, _ *mock.MockRoleRepository, _ *mock.MockMembershipRepository, _ *mock.MockSchoolRepository) {
				ur.FindByEmailFn = func(_ context.Context, _ string) (*entities.User, error) {
					return &entities.User{
						ID: uuid.New(), Email: "admin@test.com", PasswordHash: hashedPassword,
						IsActive: true,
					}, nil
				}
			},
			wantErr:   true,
			errTarget: authService.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, tokenSvc := setupAuthTestDeps()
			if tt.setupMock != nil {
				tt.setupMock(userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo)
			}

			svc := authService.NewAuthService(userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, tokenSvc, mock.NewMockLogger())
			result, err := svc.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.ErrorIs(t, err, tt.errTarget)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.Equal(t, "Bearer", result.TokenType)
				assert.NotNil(t, result.User)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "success - logout is no-op",
			token:   "some-token",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, tokenSvc := setupAuthTestDeps()
			svc := authService.NewAuthService(userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, tokenSvc, mock.NewMockLogger())
			err := svc.Logout(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantErr   bool
		errTarget error
	}{
		{
			name:      "error - refresh not implemented returns invalid refresh token",
			token:     "some-refresh-token",
			wantErr:   true,
			errTarget: authService.ErrInvalidRefreshToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, tokenSvc := setupAuthTestDeps()
			svc := authService.NewAuthService(userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, tokenSvc, mock.NewMockLogger())
			result, err := svc.RefreshToken(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.ErrorIs(t, err, tt.errTarget)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
