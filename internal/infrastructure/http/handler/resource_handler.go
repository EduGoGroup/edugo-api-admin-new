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
// @Accept json
// @Produce json
// @Success 200 {array} dto.ResourceDTO
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /resources [get]
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
// @Accept json
// @Produce json
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} dto.ResourceDTO
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /resources/{id} [get]
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
// @Accept json
// @Produce json
// @Param request body dto.CreateResourceRequest true "Resource data"
// @Success 201 {object} dto.ResourceDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /resources [post]
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
// @Accept json
// @Produce json
// @Param id path string true "Resource ID (UUID)"
// @Param request body dto.UpdateResourceRequest true "Resource update data"
// @Success 200 {object} dto.ResourceDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /resources/{id} [put]
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
