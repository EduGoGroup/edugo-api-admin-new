package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

// AcademicUnitHandler handles academic unit HTTP endpoints
type AcademicUnitHandler struct {
	unitService service.AcademicUnitService
	logger      logger.Logger
}

// NewAcademicUnitHandler creates a new AcademicUnitHandler
func NewAcademicUnitHandler(unitService service.AcademicUnitService, logger logger.Logger) *AcademicUnitHandler {
	return &AcademicUnitHandler{unitService: unitService, logger: logger}
}

// CreateUnit godoc
// @Summary Create an academic unit under a school
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Param request body dto.CreateAcademicUnitRequest true "Academic unit data"
// @Success 201 {object} dto.AcademicUnitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id}/units [post]
func (h *AcademicUnitHandler) CreateUnit(c *gin.Context) {
	schoolID := c.Param("id")
	var req dto.CreateAcademicUnitRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	unit, err := h.unitService.CreateUnit(withActor(c), schoolID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, unit)
}

// ListUnitsBySchool godoc
// @Summary List academic units by school
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Param page query int false "Page number (1-based)" minimum(1)
// @Param limit query int false "Number of items per page" minimum(1)
// @Param search query string false "Search term (ILIKE)"
// @Param search_fields query string false "Comma-separated fields to search"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id}/units [get]
func (h *AcademicUnitHandler) ListUnitsBySchool(c *gin.Context) {
	schoolID := c.Param("id")
	var filters sharedrepo.ListFilters
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "limit must be a positive integer", Code: "INVALID_REQUEST"})
			return
		}
		filters.Limit = limit
	}
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "page must be a positive integer", Code: "INVALID_REQUEST"})
			return
		}
		filters.Page = page
	}
	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			filters.SearchFields = strings.Split(fields, ",")
		}
	}
	units, total, err := h.unitService.ListUnitsBySchool(c.Request.Context(), schoolID, filters)
	if err != nil {
		_ = c.Error(err)
		return
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}
	c.JSON(http.StatusOK, dto.NewPaginatedResponse(units, total, page, filters.Limit))
}

// GetUnitTree godoc
// @Summary Get hierarchical unit tree for a school
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Success 200 {object} dto.AcademicUnitResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id}/units/tree [get]
func (h *AcademicUnitHandler) GetUnitTree(c *gin.Context) {
	schoolID := c.Param("id")
	tree, err := h.unitService.GetUnitTree(c.Request.Context(), schoolID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, tree)
}

// ListUnitsByType godoc
// @Summary List academic units by type for a school
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Param type query string false "Unit type filter"
// @Param page query int false "Page number (1-based)" minimum(1)
// @Param limit query int false "Number of items per page" minimum(1)
// @Param search query string false "Search term (ILIKE)"
// @Param search_fields query string false "Comma-separated fields to search"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id}/units/by-type [get]
func (h *AcademicUnitHandler) ListUnitsByType(c *gin.Context) {
	schoolID := c.Param("id")
	unitType := c.Query("type")
	var filters sharedrepo.ListFilters
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "limit must be a positive integer", Code: "INVALID_REQUEST"})
			return
		}
		filters.Limit = limit
	}
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "page must be a positive integer", Code: "INVALID_REQUEST"})
			return
		}
		filters.Page = page
	}
	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			filters.SearchFields = strings.Split(fields, ",")
		}
	}
	units, total, err := h.unitService.ListUnitsByType(c.Request.Context(), schoolID, unitType, filters)
	if err != nil {
		_ = c.Error(err)
		return
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}
	c.JSON(http.StatusOK, dto.NewPaginatedResponse(units, total, page, filters.Limit))
}

// GetUnit godoc
// @Summary Get an academic unit by ID
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "Academic Unit ID (UUID)"
// @Success 200 {object} dto.AcademicUnitResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /units/{id} [get]
func (h *AcademicUnitHandler) GetUnit(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.unitService.GetUnit(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, unit)
}

// UpdateUnit godoc
// @Summary Update an academic unit
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "Academic Unit ID (UUID)"
// @Param request body dto.UpdateAcademicUnitRequest true "Academic unit update data"
// @Success 200 {object} dto.AcademicUnitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /units/{id} [put]
func (h *AcademicUnitHandler) UpdateUnit(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateAcademicUnitRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	unit, err := h.unitService.UpdateUnit(withActor(c), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, unit)
}

// DeleteUnit godoc
// @Summary Soft delete an academic unit
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "Academic Unit ID (UUID)"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /units/{id} [delete]
func (h *AcademicUnitHandler) DeleteUnit(c *gin.Context) {
	id := c.Param("id")
	if err := h.unitService.DeleteUnit(withActor(c), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// RestoreUnit godoc
// @Summary Restore a soft-deleted academic unit
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "Academic Unit ID (UUID)"
// @Success 200 {object} dto.AcademicUnitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /units/{id}/restore [post]
func (h *AcademicUnitHandler) RestoreUnit(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.unitService.RestoreUnit(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, unit)
}

// GetHierarchyPath godoc
// @Summary Get hierarchy path from root to unit
// @Tags academic-units
// @Accept json
// @Produce json
// @Param id path string true "Academic Unit ID (UUID)"
// @Success 200 {array} dto.AcademicUnitResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /units/{id}/hierarchy-path [get]
func (h *AcademicUnitHandler) GetHierarchyPath(c *gin.Context) {
	id := c.Param("id")
	path, err := h.unitService.GetHierarchyPath(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, path)
}
