package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type ResourceHandler struct {
	resourceService service.ResourceService
	logger          logger.Logger
}

func NewResourceHandler(resourceService service.ResourceService, logger logger.Logger) *ResourceHandler {
	return &ResourceHandler{resourceService: resourceService, logger: logger}
}

// ListResources godoc
// @Summary List all resources
// @Tags resources
// @Router /v1/resources [get]
func (h *ResourceHandler) ListResources(c *gin.Context) {
	resources, err := h.resourceService.ListResources(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resources)
}

// GetResource godoc
// @Summary Get a resource by ID
// @Tags resources
// @Router /v1/resources/{id} [get]
func (h *ResourceHandler) GetResource(c *gin.Context) {
	id := c.Param("id")
	resource, err := h.resourceService.GetResource(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resource)
}

// CreateResource godoc
// @Summary Create a resource
// @Tags resources
// @Router /v1/resources [post]
func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	resource, err := h.resourceService.CreateResource(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resource)
}

// UpdateResource godoc
// @Summary Update a resource
// @Tags resources
// @Router /v1/resources/{id} [put]
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	resource, err := h.resourceService.UpdateResource(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resource)
}
