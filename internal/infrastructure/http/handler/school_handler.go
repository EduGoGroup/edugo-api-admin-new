package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// SchoolHandler handles school HTTP endpoints
type SchoolHandler struct {
	schoolService service.SchoolService
	logger        logger.Logger
}

// NewSchoolHandler creates a new SchoolHandler
func NewSchoolHandler(schoolService service.SchoolService, logger logger.Logger) *SchoolHandler {
	return &SchoolHandler{schoolService: schoolService, logger: logger}
}

// CreateSchool godoc
// @Summary Create a new school
// @Tags schools
// @Accept json
// @Produce json
// @Param request body dto.CreateSchoolRequest true "School data"
// @Success 201 {object} dto.SchoolResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools [post]
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
	var req dto.CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	school, err := h.schoolService.CreateSchool(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, school)
}

// GetSchool godoc
// @Summary Get a school by ID
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id} [get]
func (h *SchoolHandler) GetSchool(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}
	school, err := h.schoolService.GetSchool(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, school)
}

// GetSchoolByCode godoc
// @Summary Get a school by code
// @Tags schools
// @Accept json
// @Produce json
// @Param code path string true "School code"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/code/{code} [get]
func (h *SchoolHandler) GetSchoolByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "school code is required", Code: "INVALID_REQUEST"})
		return
	}
	school, err := h.schoolService.GetSchoolByCode(c.Request.Context(), code)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, school)
}

// ListSchools godoc
// @Summary List all schools
// @Tags schools
// @Accept json
// @Produce json
// @Success 200 {array} dto.SchoolResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools [get]
func (h *SchoolHandler) ListSchools(c *gin.Context) {
	schools, err := h.schoolService.ListSchools(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, schools)
}

// UpdateSchool godoc
// @Summary Update a school
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Param request body dto.UpdateSchoolRequest true "School update data"
// @Success 200 {object} dto.SchoolResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id} [put]
func (h *SchoolHandler) UpdateSchool(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}
	var req dto.UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	school, err := h.schoolService.UpdateSchool(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, school)
}

// DeleteSchool godoc
// @Summary Delete a school
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id} [delete]
func (h *SchoolHandler) DeleteSchool(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "school ID is required", Code: "INVALID_REQUEST"})
		return
	}
	if err := h.schoolService.DeleteSchool(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
