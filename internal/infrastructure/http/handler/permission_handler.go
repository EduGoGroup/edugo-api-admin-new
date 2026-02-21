package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type PermissionHandler struct {
	permissionService service.PermissionService
	logger            logger.Logger
}

func NewPermissionHandler(permissionService service.PermissionService, logger logger.Logger) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService, logger: logger}
}

// ListPermissions godoc
// @Summary List all permissions
// @Tags permissions
// @Router /v1/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	perms, err := h.permissionService.ListPermissions(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, perms)
}

// GetPermission godoc
// @Summary Get a permission by ID
// @Tags permissions
// @Router /v1/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")
	perm, err := h.permissionService.GetPermission(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, perm)
}
