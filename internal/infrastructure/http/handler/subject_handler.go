package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type SubjectHandler struct {
	subjectService service.SubjectService
	logger         logger.Logger
}

func NewSubjectHandler(subjectService service.SubjectService, logger logger.Logger) *SubjectHandler {
	return &SubjectHandler{subjectService: subjectService, logger: logger}
}

// CreateSubject godoc
// @Summary Create a subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param request body dto.CreateSubjectRequest true "Subject data"
// @Success 201 {object} dto.SubjectResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	subject, err := h.subjectService.CreateSubject(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, subject)
}

// ListSubjects godoc
// @Summary List all subjects
// @Tags subjects
// @Accept json
// @Produce json
// @Success 200 {array} dto.SubjectResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects [get]
func (h *SubjectHandler) ListSubjects(c *gin.Context) {
	subjects, err := h.subjectService.ListSubjects(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, subjects)
}

// GetSubject godoc
// @Summary Get a subject by ID
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID (UUID)"
// @Success 200 {object} dto.SubjectResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects/{id} [get]
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	id := c.Param("id")
	subject, err := h.subjectService.GetSubject(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, subject)
}

// UpdateSubject godoc
// @Summary Update a subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID (UUID)"
// @Param request body dto.UpdateSubjectRequest true "Subject update data"
// @Success 200 {object} dto.SubjectResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects/{id} [patch]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	subject, err := h.subjectService.UpdateSubject(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, subject)
}

// DeleteSubject godoc
// @Summary Delete a subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID (UUID)"
// @Success 204 "No content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id := c.Param("id")
	if err := h.subjectService.DeleteSubject(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
