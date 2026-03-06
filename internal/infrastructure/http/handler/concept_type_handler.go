package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// ConceptTypeHandler handles concept type HTTP endpoints
type ConceptTypeHandler struct {
	conceptTypeService service.ConceptTypeService
	logger             logger.Logger
}

// NewConceptTypeHandler creates a new ConceptTypeHandler
func NewConceptTypeHandler(conceptTypeService service.ConceptTypeService, logger logger.Logger) *ConceptTypeHandler {
	return &ConceptTypeHandler{conceptTypeService: conceptTypeService, logger: logger}
}

// CreateConceptType godoc
// @Summary Create a new concept type
// @Tags concept-types
// @Accept json
// @Produce json
// @Param request body dto.CreateConceptTypeRequest true "Concept type data"
// @Success 201 {object} dto.ConceptTypeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types [post]
func (h *ConceptTypeHandler) CreateConceptType(c *gin.Context) {
	var req dto.CreateConceptTypeRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	ct, err := h.conceptTypeService.CreateConceptType(withActor(c), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, ct)
}

// ListConceptTypes godoc
// @Summary List all concept types
// @Tags concept-types
// @Accept json
// @Produce json
// @Success 200 {array} dto.ConceptTypeResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types [get]
func (h *ConceptTypeHandler) ListConceptTypes(c *gin.Context) {
	types, err := h.conceptTypeService.ListConceptTypes(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, types)
}

// GetConceptType godoc
// @Summary Get a concept type by ID
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Success 200 {object} dto.ConceptTypeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id} [get]
func (h *ConceptTypeHandler) GetConceptType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	ct, err := h.conceptTypeService.GetConceptType(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ct)
}

// UpdateConceptType godoc
// @Summary Update a concept type
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Param request body dto.UpdateConceptTypeRequest true "Concept type update data"
// @Success 200 {object} dto.ConceptTypeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id} [put]
func (h *ConceptTypeHandler) UpdateConceptType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	var req dto.UpdateConceptTypeRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	ct, err := h.conceptTypeService.UpdateConceptType(withActor(c), id, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ct)
}

// DeleteConceptType godoc
// @Summary Delete a concept type
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id} [delete]
func (h *ConceptTypeHandler) DeleteConceptType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	if err := h.conceptTypeService.DeleteConceptType(withActor(c), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateDefinition godoc
// @Summary Create a concept definition
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Param request body dto.ConceptDefinitionRequest true "Definition data"
// @Success 201 {object} dto.ConceptDefinitionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id}/definitions [post]
func (h *ConceptTypeHandler) CreateDefinition(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	var req dto.ConceptDefinitionRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	def, err := h.conceptTypeService.CreateDefinition(withActor(c), typeID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, def)
}

// ListDefinitions godoc
// @Summary List concept definitions by type
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Success 200 {array} dto.ConceptDefinitionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id}/definitions [get]
func (h *ConceptTypeHandler) ListDefinitions(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	defs, err := h.conceptTypeService.ListDefinitions(c.Request.Context(), typeID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, defs)
}

// UpdateDefinition godoc
// @Summary Update a concept definition
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Param defId path string true "Definition ID (UUID)"
// @Param request body dto.ConceptDefinitionRequest true "Definition update data"
// @Success 200 {object} dto.ConceptDefinitionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id}/definitions/{defId} [put]
func (h *ConceptTypeHandler) UpdateDefinition(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	defID, err := uuid.Parse(c.Param("defId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid definition ID", Code: "INVALID_REQUEST"})
		return
	}
	var req dto.ConceptDefinitionRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	def, err := h.conceptTypeService.UpdateDefinition(withActor(c), typeID, defID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, def)
}

// DeleteDefinition godoc
// @Summary Delete a concept definition
// @Tags concept-types
// @Accept json
// @Produce json
// @Param id path string true "Concept Type ID (UUID)"
// @Param defId path string true "Definition ID (UUID)"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /concept-types/{id}/definitions/{defId} [delete]
func (h *ConceptTypeHandler) DeleteDefinition(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept type ID", Code: "INVALID_REQUEST"})
		return
	}
	defID, err := uuid.Parse(c.Param("defId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid definition ID", Code: "INVALID_REQUEST"})
		return
	}
	if err := h.conceptTypeService.DeleteDefinition(withActor(c), typeID, defID); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetSchoolConcepts godoc
// @Summary Get school concepts
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Success 200 {array} dto.SchoolConceptResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id}/concepts [get]
func (h *ConceptTypeHandler) GetSchoolConcepts(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid school ID", Code: "INVALID_REQUEST"})
		return
	}
	concepts, err := h.conceptTypeService.GetSchoolConcepts(c.Request.Context(), schoolID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, concepts)
}

// UpdateSchoolConcept godoc
// @Summary Update a school concept
// @Tags schools
// @Accept json
// @Produce json
// @Param id path string true "School ID (UUID)"
// @Param conceptId path string true "Concept ID (UUID)"
// @Param request body dto.UpdateSchoolConceptRequest true "School concept update data"
// @Success 200 {object} dto.SchoolConceptResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /schools/{id}/concepts/{conceptId} [put]
func (h *ConceptTypeHandler) UpdateSchoolConcept(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid school ID", Code: "INVALID_REQUEST"})
		return
	}
	conceptID, err := uuid.Parse(c.Param("conceptId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid concept ID", Code: "INVALID_REQUEST"})
		return
	}
	var req dto.UpdateSchoolConceptRequest
	if err := bindJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	concept, err := h.conceptTypeService.UpdateSchoolConcept(withActor(c), schoolID, conceptID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, concept)
}
