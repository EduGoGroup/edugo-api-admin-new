package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// Sentinel errors for auth operations
var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserInactive        = errors.New("user inactive")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

// AuthService defines the authentication service interface
type AuthService interface {
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}

type authService struct {
	userRepo     repository.UserRepository
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
	tokenService *TokenService
	logger       logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	tokenService *TokenService,
	logger logger.Logger,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		tokenService: tokenService,
		logger:       logger,
	}
}

// Login validates credentials and returns JWT tokens
func (s *authService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		s.logger.Warn("login attempt with non-existent email", "email", email)
		return nil, ErrInvalidCredentials
	}

	// 2. Verify user is active
	if !user.IsActive {
		s.logger.Warn("login attempt with inactive user", "email", email, "user_id", user.ID.String())
		return nil, ErrUserInactive
	}

	// 3. Verify password
	if err := auth.VerifyPassword(user.PasswordHash, password); err != nil {
		s.logger.Warn("incorrect password", "email", email)
		return nil, ErrInvalidCredentials
	}

	// 4. Get school_id
	schoolID := ""
	if user.SchoolID != nil {
		schoolID = user.SchoolID.String()
	}

	// 5. Build RBAC context
	activeContext := s.buildUserContext(ctx, user.ID, user.SchoolID)
	if activeContext == nil {
		s.logger.Error("no RBAC context found for user", "user_id", user.ID.String(), "email", user.Email)
		return nil, fmt.Errorf("user has no assigned roles")
	}

	// 6. Generate tokens
	tokenResponse, err := s.tokenService.GenerateTokenPairWithContext(user.ID.String(), user.Email, activeContext)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// 7. Add user info
	tokenResponse.User = &dto.UserInfo{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FirstName + " " + user.LastName,
		SchoolID:  schoolID,
	}

	tokenResponse.ActiveContext = &dto.UserContextDTO{
		RoleID:      activeContext.RoleID,
		RoleName:    activeContext.RoleName,
		SchoolID:    activeContext.SchoolID,
		Permissions: activeContext.Permissions,
	}

	s.logger.Info("user logged in",
		"entity_type", "auth_session",
		"user_id", user.ID.String(),
		"email", user.Email,
		"role", activeContext.RoleName,
		"school_id", schoolID,
	)

	// 8. Update last login (fire and forget)
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		user.UpdatedAt = time.Now()
		if err := s.userRepo.Update(bgCtx, user); err != nil {
			s.logger.Warn("error updating last login", "error", err)
		}
	}()

	return tokenResponse, nil
}

// Logout invalidates the access token
func (s *authService) Logout(_ context.Context, _ string) error {
	// Token blacklisting requires Redis (not implemented in this clean API)
	// For now, logout is a no-op on the server side; the client discards the token
	s.logger.Info("user logged out", "entity_type", "auth_session")
	return nil
}

// RefreshToken validates a refresh token and generates a new access token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// Verify the refresh token by checking its hash
	tokenHash := auth.HashToken(refreshToken)
	if tokenHash == "" {
		return nil, ErrInvalidRefreshToken
	}

	// Since refresh tokens are opaque (not JWT), we validate by attempting
	// to extract user info from the stored token. In this clean implementation,
	// refresh tokens are stateless secure random tokens.
	// A production implementation would store refresh token hashes in the DB.
	// For now, we validate the format and rely on the access token flow.
	_ = tokenHash

	// Without a refresh token store, we cannot validate refresh tokens properly.
	// This is a known limitation that requires Redis or DB storage.
	return nil, ErrInvalidRefreshToken
}

// buildUserContext constructs the RBAC UserContext for the JWT
func (s *authService) buildUserContext(ctx context.Context, userID uuid.UUID, schoolID *uuid.UUID) *auth.UserContext {
	// Get user roles in this context
	userRoles, err := s.userRoleRepo.FindByUserInContext(ctx, userID, schoolID, nil)
	if err != nil {
		s.logger.Warn("error obtaining user roles for RBAC context",
			"user_id", userID.String(),
			"error", err,
		)
		return nil
	}
	if len(userRoles) == 0 {
		return nil
	}

	// Use the first active role as the primary context
	firstRole := userRoles[0]
	role, err := s.roleRepo.FindByID(ctx, firstRole.RoleID)
	if err != nil {
		s.logger.Warn("error obtaining role for RBAC context",
			"user_id", userID.String(),
			"role_id", firstRole.RoleID.String(),
			"error", err,
		)
		return nil
	}

	// Get user permissions in this context
	permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID, schoolID, nil)
	if err != nil {
		s.logger.Warn("error obtaining user permissions",
			"user_id", userID.String(),
			"error", err,
		)
		permissions = []string{}
	}

	uc := &auth.UserContext{
		RoleID:      role.ID.String(),
		RoleName:    role.Name,
		Permissions: permissions,
	}

	if schoolID != nil {
		uc.SchoolID = schoolID.String()
	}

	return uc
}
