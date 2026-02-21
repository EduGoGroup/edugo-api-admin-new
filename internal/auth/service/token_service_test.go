package service_test

import (
	"context"
	"testing"
	"time"

	authService "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/service"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_GenerateTokenPairWithContext(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key-for-unit-tests", "edugo-test")

	tests := []struct {
		name      string
		userID    string
		email     string
		context   *auth.UserContext
		wantErr   bool
	}{
		{
			name:   "success - generates token pair",
			userID: "user-123",
			email:  "user@test.com",
			context: &auth.UserContext{
				RoleID:      "role-1",
				RoleName:    "admin",
				Permissions: []string{"schools:read", "schools:create"},
			},
			wantErr: false,
		},
		{
			name:   "success - with minimal context",
			userID: "user-456",
			email:  "other@test.com",
			context: &auth.UserContext{
				RoleID:      "role-2",
				RoleName:    "viewer",
				Permissions: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := authService.NewTokenService(jwtManager, 15*time.Minute, 7*24*time.Hour)
			result, err := svc.GenerateTokenPairWithContext(tt.userID, tt.email, tt.context)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.Equal(t, "Bearer", result.TokenType)
				assert.Greater(t, result.ExpiresIn, int64(0))
			}
		})
	}
}

func TestTokenService_VerifyToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key-for-unit-tests", "edugo-test")

	validToken, _, _ := jwtManager.GenerateTokenWithContext("user-123", "user@test.com", &auth.UserContext{
		RoleID:   "role-1",
		RoleName: "admin",
	}, 15*time.Minute)

	tests := []struct {
		name      string
		token     string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "success - valid token",
			token:     validToken,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "invalid token returns valid=false",
			token:     "invalid-jwt-token",
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := authService.NewTokenService(jwtManager, 15*time.Minute, 7*24*time.Hour)
			result, err := svc.VerifyToken(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.wantValid, result.Valid)
				if tt.wantValid {
					assert.Equal(t, "user-123", result.UserID)
					assert.Equal(t, "user@test.com", result.Email)
				}
			}
		})
	}
}

func TestTokenService_DefaultDurations(t *testing.T) {
	jwtManager := auth.NewJWTManager("secret", "issuer")

	t.Run("uses default durations when zero", func(t *testing.T) {
		svc := authService.NewTokenService(jwtManager, 0, 0)
		result, err := svc.GenerateTokenPairWithContext("user-1", "a@b.com", &auth.UserContext{
			RoleID: "r1", RoleName: "admin", Permissions: []string{},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
	})
}
