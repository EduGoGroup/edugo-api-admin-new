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

func (h *ScreenConfigHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.screenService.GetTemplate(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, template)
}

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

func (h *ScreenConfigHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.screenService.DeleteTemplate(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Instances

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

func (h *ScreenConfigHandler) GetInstance(c *gin.Context) {
	id := c.Param("id")
	instance, err := h.screenService.GetInstance(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, instance)
}

func (h *ScreenConfigHandler) GetInstanceByKey(c *gin.Context) {
	key := c.Param("key")
	instance, err := h.screenService.GetInstanceByKey(c.Request.Context(), key)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, instance)
}

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

func (h *ScreenConfigHandler) DeleteInstance(c *gin.Context) {
	id := c.Param("id")
	if err := h.screenService.DeleteInstance(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Resolve

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

func (h *ScreenConfigHandler) GetScreensForResource(c *gin.Context) {
	resourceID := c.Param("resourceId")
	screens, err := h.screenService.GetScreensForResource(c.Request.Context(), resourceID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, screens)
}

func (h *ScreenConfigHandler) UnlinkScreen(c *gin.Context) {
	id := c.Param("id")
	if err := h.screenService.UnlinkScreen(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
