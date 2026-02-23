package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type GuardianHandler struct {
	guardianService service.GuardianService
	logger          logger.Logger
}

func NewGuardianHandler(guardianService service.GuardianService, logger logger.Logger) *GuardianHandler {
	return &GuardianHandler{guardianService: guardianService, logger: logger}
}

// CreateGuardianRelation godoc
// @Summary Create a guardian relation
// @Tags guardian-relations
// @Accept json
// @Produce json
// @Param request body dto.CreateGuardianRelationRequest true "Guardian relation data"
// @Success 201 {object} dto.GuardianRelationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /guardian-relations [post]
func (h *GuardianHandler) CreateGuardianRelation(c *gin.Context) {
	var req dto.CreateGuardianRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	createdBy, _ := c.Get("user_id")
	createdByStr := ""
	if createdBy != nil {
		createdByStr, _ = createdBy.(string)
	}
	relation, err := h.guardianService.CreateRelation(c.Request.Context(), req, createdByStr)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, relation)
}

// GetGuardianRelation godoc
// @Summary Get a guardian relation by ID
// @Tags guardian-relations
// @Accept json
// @Produce json
// @Param id path string true "Guardian Relation ID (UUID)"
// @Success 200 {object} dto.GuardianRelationResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /guardian-relations/{id} [get]
func (h *GuardianHandler) GetGuardianRelation(c *gin.Context) {
	id := c.Param("id")
	relation, err := h.guardianService.GetRelation(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, relation)
}

// UpdateGuardianRelation godoc
// @Summary Update a guardian relation
// @Tags guardian-relations
// @Accept json
// @Produce json
// @Param id path string true "Guardian Relation ID (UUID)"
// @Param request body dto.UpdateGuardianRelationRequest true "Guardian relation update data"
// @Success 200 {object} dto.GuardianRelationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /guardian-relations/{id} [put]
func (h *GuardianHandler) UpdateGuardianRelation(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateGuardianRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	relation, err := h.guardianService.UpdateRelation(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, relation)
}

// DeleteGuardianRelation godoc
// @Summary Delete a guardian relation
// @Tags guardian-relations
// @Accept json
// @Produce json
// @Param id path string true "Guardian Relation ID (UUID)"
// @Success 204 "No content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /guardian-relations/{id} [delete]
func (h *GuardianHandler) DeleteGuardianRelation(c *gin.Context) {
	id := c.Param("id")
	if err := h.guardianService.DeleteRelation(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetGuardianRelations godoc
// @Summary Get relations for a guardian
// @Tags guardian-relations
// @Accept json
// @Produce json
// @Param guardian_id path string true "Guardian ID (UUID)"
// @Success 200 {array} dto.GuardianRelationResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /guardians/{guardian_id}/relations [get]
func (h *GuardianHandler) GetGuardianRelations(c *gin.Context) {
	guardianID := c.Param("guardian_id")
	relations, err := h.guardianService.GetGuardianRelations(c.Request.Context(), guardianID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, relations)
}

// GetStudentGuardians godoc
// @Summary Get guardians for a student
// @Tags guardian-relations
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID (UUID)"
// @Success 200 {array} dto.GuardianRelationResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /students/{student_id}/guardians [get]
func (h *GuardianHandler) GetStudentGuardians(c *gin.Context) {
	studentID := c.Param("student_id")
	relations, err := h.guardianService.GetStudentGuardians(c.Request.Context(), studentID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, relations)
}
