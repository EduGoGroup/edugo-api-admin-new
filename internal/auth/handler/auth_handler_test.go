package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	authDTO "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/dto"
	authHandler "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/handler"
	authService "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/service"
	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMock  func(m *mock.MockAuthService)
		wantStatus int
	}{
		{
			name: "success - returns 200",
			body: authDTO.LoginRequest{Email: "user@example.com", Password: "password123"},
			setupMock: func(m *mock.MockAuthService) {
				m.LoginFn = func(_ context.Context, _, _ string) (*authDTO.LoginResponse, error) {
					return &authDTO.LoginResponse{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						ExpiresIn:    900,
						TokenType:    "Bearer",
						User:         &authDTO.UserInfo{ID: "1", Email: "user@example.com"},
					}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - invalid body returns 400",
			body:       "bad",
			setupMock:  func(_ *mock.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - invalid credentials returns 401",
			body: authDTO.LoginRequest{Email: "user@example.com", Password: "wrongpassword123"},
			setupMock: func(m *mock.MockAuthService) {
				m.LoginFn = func(_ context.Context, _, _ string) (*authDTO.LoginResponse, error) {
					return nil, authService.ErrInvalidCredentials
				}
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "error - user inactive returns 403",
			body: authDTO.LoginRequest{Email: "user@example.com", Password: "password123"},
			setupMock: func(m *mock.MockAuthService) {
				m.LoginFn = func(_ context.Context, _, _ string) (*authDTO.LoginResponse, error) {
					return nil, authService.ErrUserInactive
				}
			},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := authHandler.NewAuthHandler(mockSvc, mock.NewMockLogger())
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.POST("/auth/login", h.Login)

			req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMock  func(m *mock.MockAuthService)
		wantStatus int
	}{
		{
			name:       "error - invalid body returns 400",
			body:       "bad",
			setupMock:  func(_ *mock.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - invalid refresh token returns 401",
			body: authDTO.RefreshTokenRequest{RefreshToken: "expired-token"},
			setupMock: func(m *mock.MockAuthService) {
				m.RefreshTokenFn = func(_ context.Context, _ string) (*authDTO.RefreshResponse, error) {
					return nil, authService.ErrInvalidRefreshToken
				}
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := authHandler.NewAuthHandler(mockSvc, mock.NewMockLogger())
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.POST("/auth/refresh", h.Refresh)

			req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		setupMock  func(m *mock.MockAuthService)
		wantStatus int
	}{
		{
			name:       "success",
			authHeader: "Bearer valid-token",
			setupMock: func(m *mock.MockAuthService) {
				m.LogoutFn = func(_ context.Context, _ string) error { return nil }
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - no token",
			authHeader: "",
			setupMock:  func(_ *mock.MockAuthService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := authHandler.NewAuthHandler(mockSvc, mock.NewMockLogger())
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.POST("/auth/logout", h.Logout)

			req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
