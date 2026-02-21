package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	authDTO "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/dto"
	authHandler "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/handler"
	authService "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/service"
	"github.com/EduGoGroup/edugo-shared/auth"
)

func TestVerifyHandler_VerifyToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key-for-testing-purposes", "edugo-test")
	tokenSvc := authService.NewTokenService(jwtManager, 15*time.Minute, 7*24*time.Hour)

	// Generate a valid token for testing
	validToken, _, _ := jwtManager.GenerateTokenWithContext("user-1", "user@test.com", &auth.UserContext{
		RoleID:      "role-1",
		RoleName:    "admin",
		Permissions: []string{"schools:read"},
	}, 15*time.Minute)

	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
	}{
		{
			name:       "success - valid token",
			body:       authDTO.VerifyTokenRequest{Token: validToken},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - invalid body",
			body:       "bad",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error - empty token",
			body:       authDTO.VerifyTokenRequest{Token: "   "},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "success - returns invalid token response (not error)",
			body:       authDTO.VerifyTokenRequest{Token: "invalid-jwt-token"},
			wantStatus: http.StatusOK, // verify returns valid=false in body, not an HTTP error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := authHandler.NewVerifyHandler(tokenSvc)
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.POST("/auth/verify", h.VerifyToken)

			req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
