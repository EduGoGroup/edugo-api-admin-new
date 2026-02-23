package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type ScreenConfigHandler struct {
	screenService service.ScreenConfigService
	logger        logger.Logger
}

func NewScreenConfigHandler(screenService service.ScreenConfigService, logger logger.Logger) *ScreenConfigHandler {
	return &ScreenConfigHandler{screenService: screenService, logger: logger}
}

// Templates

// CreateTemplate godoc
// @Summary Create a screen template
// @Tags screen-config
// @Accept json
// @Produce json
// @Param request body service.CreateTemplateRequest true "Template data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/templates [post]
func (h *ScreenConfigHandler) CreateTemplate(c *gin.Context) {
	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	template, err := h.screenService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, template)
}

// ListTemplates godoc
// @Summary List screen templates
// @Tags screen-config
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/templates [get]
func (h *ScreenConfigHandler) ListTemplates(c *gin.Context) {
	var filter service.TemplateFilter
	_ = c.ShouldBindQuery(&filter)
	templates, total, err := h.screenService.ListTemplates(c.Request.Context(), filter)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": templates, "total": total})
}

// GetTemplate godoc
// @Summary Get a screen template by ID
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Template ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/templates/{id} [get]
func (h *ScreenConfigHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.screenService.GetTemplate(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, template)
}

// UpdateTemplate godoc
// @Summary Update a screen template
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Template ID (UUID)"
// @Param request body service.UpdateTemplateRequest true "Template update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/templates/{id} [put]
func (h *ScreenConfigHandler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	template, err := h.screenService.UpdateTemplate(c.Request.Context(), id, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, template)
}

// DeleteTemplate godoc
// @Summary Delete a screen template
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Template ID (UUID)"
// @Success 204 "No content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/templates/{id} [delete]
func (h *ScreenConfigHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.screenService.DeleteTemplate(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Instances

// CreateInstance godoc
// @Summary Create a screen instance
// @Tags screen-config
// @Accept json
// @Produce json
// @Param request body service.CreateInstanceRequest true "Instance data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/instances [post]
func (h *ScreenConfigHandler) CreateInstance(c *gin.Context) {
	var req service.CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	instance, err := h.screenService.CreateInstance(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, instance)
}

// ListInstances godoc
// @Summary List screen instances
// @Tags screen-config
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/instances [get]
func (h *ScreenConfigHandler) ListInstances(c *gin.Context) {
	var filter service.InstanceFilter
	_ = c.ShouldBindQuery(&filter)
	instances, total, err := h.screenService.ListInstances(c.Request.Context(), filter)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instances, "total": total})
}

// GetInstance godoc
// @Summary Get a screen instance by ID
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/instances/{id} [get]
func (h *ScreenConfigHandler) GetInstance(c *gin.Context) {
	id := c.Param("id")
	instance, err := h.screenService.GetInstance(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, instance)
}

// GetInstanceByKey godoc
// @Summary Get a screen instance by key
// @Tags screen-config
// @Accept json
// @Produce json
// @Param key path string true "Instance key"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/instances/key/{key} [get]
func (h *ScreenConfigHandler) GetInstanceByKey(c *gin.Context) {
	key := c.Param("key")
	instance, err := h.screenService.GetInstanceByKey(c.Request.Context(), key)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, instance)
}

// UpdateInstance godoc
// @Summary Update a screen instance
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Param request body service.UpdateInstanceRequest true "Instance update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/instances/{id} [put]
func (h *ScreenConfigHandler) UpdateInstance(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	instance, err := h.screenService.UpdateInstance(c.Request.Context(), id, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, instance)
}

// DeleteInstance godoc
// @Summary Delete a screen instance
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Instance ID (UUID)"
// @Success 204 "No content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/instances/{id} [delete]
func (h *ScreenConfigHandler) DeleteInstance(c *gin.Context) {
	id := c.Param("id")
	if err := h.screenService.DeleteInstance(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Resolve

// ResolveScreenByKey godoc
// @Summary Resolve a screen configuration by key
// @Tags screen-config
// @Accept json
// @Produce json
// @Param key path string true "Screen key"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/resolve/key/{key} [get]
func (h *ScreenConfigHandler) ResolveScreenByKey(c *gin.Context) {
	key := c.Param("key")
	combined, err := h.screenService.ResolveScreenByKey(c.Request.Context(), key)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, combined)
}

// Resource-Screens

// LinkScreenToResource godoc
// @Summary Link a screen instance to a resource
// @Tags screen-config
// @Accept json
// @Produce json
// @Param request body service.LinkScreenRequest true "Link data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/resource-screens [post]
func (h *ScreenConfigHandler) LinkScreenToResource(c *gin.Context) {
	var req service.LinkScreenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	rs, err := h.screenService.LinkScreenToResource(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, rs)
}

// GetScreensForResource godoc
// @Summary Get screens linked to a resource
// @Tags screen-config
// @Accept json
// @Produce json
// @Param resourceId path string true "Resource ID (UUID)"
// @Success 200 {array} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/resource-screens/{resourceId} [get]
func (h *ScreenConfigHandler) GetScreensForResource(c *gin.Context) {
	resourceID := c.Param("resourceId")
	screens, err := h.screenService.GetScreensForResource(c.Request.Context(), resourceID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, screens)
}

// UnlinkScreen godoc
// @Summary Unlink a screen from a resource
// @Tags screen-config
// @Accept json
// @Produce json
// @Param id path string true "Resource-Screen link ID (UUID)"
// @Success 204 "No content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /screen-config/resource-screens/{id} [delete]
func (h *ScreenConfigHandler) UnlinkScreen(c *gin.Context) {
	id := c.Param("id")
	if err := h.screenService.UnlinkScreen(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
