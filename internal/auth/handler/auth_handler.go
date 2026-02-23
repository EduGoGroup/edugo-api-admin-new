package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/auth/service"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthService
	logger      logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: log}
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Email and password are required",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid credentials",
				Code:    "INVALID_CREDENTIALS",
			})
		case errors.Is(err, service.ErrUserInactive):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "User inactive",
				Code:    "USER_INACTIVE",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Authentication error",
				Code:    "AUTH_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generates a new access token using the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.RefreshResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Refresh token is required",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired refresh token",
				Code:    "INVALID_REFRESH_TOKEN",
			})
		case errors.Is(err, service.ErrUserInactive):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "User inactive",
				Code:    "USER_INACTIVE",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Error refreshing tokens",
				Code:    "REFRESH_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary User logout
// @Description Invalidates the current access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Authorization token required",
			Code:    "TOKEN_REQUIRED",
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Logout error",
			Code:    "LOGOUT_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// SwitchContext godoc
// @Summary Switch school context
// @Description Switches the active school context for the authenticated user and returns new tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SwitchContextRequest true "Target school"
// @Success 200 {object} dto.SwitchContextResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /auth/switch-context [post]
func (h *AuthHandler) SwitchContext(c *gin.Context) {
	userID, err := ginmiddleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    "NOT_AUTHENTICATED",
		})
		return
	}

	var req dto.SwitchContextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "school_id is required and must be a valid UUID",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	response, err := h.authService.SwitchContext(c.Request.Context(), userID, req.SchoolID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoMembership):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "No active membership in target school",
				Code:    "NO_MEMBERSHIP",
			})
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "User not found",
				Code:    "USER_NOT_FOUND",
			})
		case errors.Is(err, service.ErrUserInactive):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "User inactive",
				Code:    "USER_INACTIVE",
			})
		case errors.Is(err, service.ErrInvalidSchoolID):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "bad_request",
				Message: "Invalid school_id",
				Code:    "INVALID_SCHOOL_ID",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Error switching context",
				Code:    "SWITCH_CONTEXT_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAvailableContexts godoc
// @Summary Get available contexts
// @Description Returns all available contexts (roles/schools) for the authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} dto.AvailableContextsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /auth/contexts [get]
func (h *AuthHandler) GetAvailableContexts(c *gin.Context) {
	userID, err := ginmiddleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    "NOT_AUTHENTICATED",
		})
		return
	}

	// Get current context from JWT claims
	claims, _ := ginmiddleware.GetClaims(c)
	var currentContext *auth.UserContext
	if claims != nil {
		currentContext = claims.ActiveContext
	}

	response, err := h.authService.GetAvailableContexts(c.Request.Context(), userID, currentContext)
	if err != nil {
		h.logger.Error("error fetching available contexts", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Error fetching available contexts",
			Code:    "CONTEXTS_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
