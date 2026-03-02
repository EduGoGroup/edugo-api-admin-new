package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/middleware"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

type SubjectHandler struct {
	subjectService service.SubjectService
	logger         logger.Logger
}

func NewSubjectHandler(subjectService service.SubjectService, logger logger.Logger) *SubjectHandler {
	return &SubjectHandler{subjectService: subjectService, logger: logger}
}

// extractSchoolID extracts the school ID from the JWT active context
func (h *SubjectHandler) extractSchoolID(c *gin.Context) (string, bool) {
	val, exists := c.Get(middleware.ContextKeyActiveContext)
	if !exists {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "no active context", Code: "NO_ACTIVE_CONTEXT"})
		return "", false
	}
	ac, ok := val.(*auth.UserContext)
	if !ok || ac.SchoolID == "" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "no school context", Code: "NO_SCHOOL_CONTEXT"})
		return "", false
	}
	return ac.SchoolID, true
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
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	schoolID, ok := h.extractSchoolID(c)
	if !ok {
		return
	}
	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	subject, err := h.subjectService.CreateSubject(c.Request.Context(), schoolID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, subject)
}

// ListSubjects godoc
// @Summary List subjects for the current school
// @Tags subjects
// @Accept json
// @Produce json
// @Param search query string false "Search term (ILIKE)"
// @Param search_fields query string false "Comma-separated fields to search"
// @Success 200 {array} dto.SubjectResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /subjects [get]
func (h *SubjectHandler) ListSubjects(c *gin.Context) {
	schoolID, ok := h.extractSchoolID(c)
	if !ok {
		return
	}
	var filters sharedrepo.ListFilters
	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			filters.SearchFields = strings.Split(fields, ",")
		}
	}
	subjects, err := h.subjectService.ListSubjects(c.Request.Context(), schoolID, filters)
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
