package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// MaterialHandler handles material HTTP endpoints
type MaterialHandler struct {
	materialService service.MaterialService
	logger          logger.Logger
}

// NewMaterialHandler creates a new MaterialHandler
func NewMaterialHandler(materialService service.MaterialService, logger logger.Logger) *MaterialHandler {
	return &MaterialHandler{materialService: materialService, logger: logger}
}

// DeleteMaterial godoc
// @Summary Delete a material
// @Tags materials
// @Accept json
// @Produce json
// @Param id path string true "Material ID (UUID)"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /materials/{id} [delete]
func (h *MaterialHandler) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "material ID is required", Code: "INVALID_REQUEST"})
		return
	}
	if err := h.materialService.DeleteMaterial(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
