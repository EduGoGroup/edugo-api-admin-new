package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

// UserHandler handles user HTTP endpoints
type UserHandler struct {
	userService service.UserService
	logger      logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// ListUsers godoc
// @Summary List all users
// @Tags users
// @Accept json
// @Produce json
// @Param is_active query bool false "Filter by active status"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var filters sharedrepo.ListFilters

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid is_active parameter", Code: "INVALID_REQUEST"})
			return
		}
		filters.IsActive = &isActive
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid limit parameter", Code: "INVALID_REQUEST"})
			return
		}
		filters.Limit = limit
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid offset parameter", Code: "INVALID_REQUEST"})
			return
		}
		filters.Offset = offset
	}

	users, err := h.userService.ListUsers(c.Request.Context(), filters)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "user ID is required", Code: "INVALID_REQUEST"})
		return
	}
	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UpdateUserRequest true "User update data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "user ID is required", Code: "INVALID_REQUEST"})
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "user ID is required", Code: "INVALID_REQUEST"})
		return
	}
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
