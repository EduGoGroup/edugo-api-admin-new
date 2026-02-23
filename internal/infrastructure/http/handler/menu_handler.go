package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type MenuHandler struct {
	menuService service.MenuService
	logger      logger.Logger
}

func NewMenuHandler(menuService service.MenuService, logger logger.Logger) *MenuHandler {
	return &MenuHandler{menuService: menuService, logger: logger}
}

// GetUserMenu godoc
// @Summary Get menu filtered by user permissions
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /menu [get]
func (h *MenuHandler) GetUserMenu(c *gin.Context) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	jwtClaims, ok := claims.(*auth.Claims)
	if !ok || jwtClaims.ActiveContext == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no active context"})
		return
	}

	menu, err := h.menuService.GetMenuForUser(c.Request.Context(), jwtClaims.ActiveContext.Permissions)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, menu)
}

// GetFullMenu godoc
// @Summary Get full menu (admin view)
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /menu/full [get]
func (h *MenuHandler) GetFullMenu(c *gin.Context) {
	menu, err := h.menuService.GetFullMenu(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, menu)
}
