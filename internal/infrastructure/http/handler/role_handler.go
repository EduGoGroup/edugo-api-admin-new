package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type RoleHandler struct {
	roleService service.RoleService
	logger      logger.Logger
}

func NewRoleHandler(roleService service.RoleService, logger logger.Logger) *RoleHandler {
	return &RoleHandler{roleService: roleService, logger: logger}
}

// ListRoles godoc
// @Summary List all roles
// @Tags roles
// @Router /v1/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	scope := c.Query("scope")
	roles, err := h.roleService.GetRoles(c.Request.Context(), scope)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GetRole godoc
// @Summary Get a role by ID
// @Tags roles
// @Router /v1/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	role, err := h.roleService.GetRole(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, role)
}

// GetRolePermissions godoc
// @Summary Get permissions for a role
// @Tags roles
// @Router /v1/roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id := c.Param("id")
	perms, err := h.roleService.GetRolePermissions(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, perms)
}

// GetUserRoles godoc
// @Summary Get roles assigned to a user
// @Tags users
// @Router /v1/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")
	roles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GrantRole godoc
// @Summary Grant a role to a user
// @Tags users
// @Router /v1/users/{user_id}/roles [post]
func (h *RoleHandler) GrantRole(c *gin.Context) {
	userID := c.Param("user_id")
	var req dto.GrantRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	grantedBy, _ := c.Get("user_id")
	grantedByStr := ""
	if grantedBy != nil {
		grantedByStr, _ = grantedBy.(string)
	}
	result, err := h.roleService.GrantRoleToUser(c.Request.Context(), userID, &req, grantedByStr)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// RevokeRole godoc
// @Summary Revoke a role from a user
// @Tags users
// @Router /v1/users/{user_id}/roles/{role_id} [delete]
func (h *RoleHandler) RevokeRole(c *gin.Context) {
	userID := c.Param("user_id")
	roleID := c.Param("role_id")
	if err := h.roleService.RevokeRoleFromUser(c.Request.Context(), userID, roleID); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
