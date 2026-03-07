package dto

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateConceptTypeRequest represents the request to create a concept type
type CreateConceptTypeRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Code        string `json:"code" binding:"required,min=3"`
	Description string `json:"description"`
}

// UpdateConceptTypeRequest represents the request to update a concept type
type UpdateConceptTypeRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// ConceptTypeResponse represents a concept type in API responses
type ConceptTypeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// ConceptDefinitionRequest represents the request to create/update a concept definition
type ConceptDefinitionRequest struct {
	TermKey   string `json:"term_key" binding:"required"`
	TermValue string `json:"term_value" binding:"required"`
	Category  string `json:"category"`
	SortOrder int    `json:"sort_order"`
}

// ConceptDefinitionResponse represents a concept definition in API responses
type ConceptDefinitionResponse struct {
	ID        string `json:"id"`
	TermKey   string `json:"term_key"`
	TermValue string `json:"term_value"`
	Category  string `json:"category"`
	SortOrder int    `json:"sort_order"`
}

// SchoolConceptResponse represents a school concept in API responses
type SchoolConceptResponse struct {
	ID        string `json:"id"`
	TermKey   string `json:"term_key"`
	TermValue string `json:"term_value"`
	Category  string `json:"category"`
}

// UpdateSchoolConceptRequest represents the request to update a school concept
type UpdateSchoolConceptRequest struct {
	TermValue string `json:"term_value" binding:"required"`
}

// ToConceptTypeResponse converts a ConceptType entity to ConceptTypeResponse
func ToConceptTypeResponse(ct *entities.ConceptType) ConceptTypeResponse {
	var description string
	if ct.Description != nil {
		description = *ct.Description
	}
	return ConceptTypeResponse{
		ID:          ct.ID.String(),
		Name:        ct.Name,
		Code:        ct.Code,
		Description: description,
		IsActive:    ct.IsActive,
	}
}

// ToConceptTypeResponseList converts a slice of ConceptType entities to ConceptTypeResponse slice
func ToConceptTypeResponseList(types []*entities.ConceptType) []ConceptTypeResponse {
	responses := make([]ConceptTypeResponse, len(types))
	for i, ct := range types {
		responses[i] = ToConceptTypeResponse(ct)
	}
	return responses
}

// ToConceptDefinitionResponse converts a ConceptDefinition entity to ConceptDefinitionResponse
func ToConceptDefinitionResponse(def *entities.ConceptDefinition) ConceptDefinitionResponse {
	return ConceptDefinitionResponse{
		ID:        def.ID.String(),
		TermKey:   def.TermKey,
		TermValue: def.TermValue,
		Category:  def.Category,
		SortOrder: def.SortOrder,
	}
}

// ToConceptDefinitionResponseList converts a slice of ConceptDefinition entities to ConceptDefinitionResponse slice
func ToConceptDefinitionResponseList(defs []*entities.ConceptDefinition) []ConceptDefinitionResponse {
	responses := make([]ConceptDefinitionResponse, len(defs))
	for i, def := range defs {
		responses[i] = ToConceptDefinitionResponse(def)
	}
	return responses
}

// ToSchoolConceptResponse converts a SchoolConcept entity to SchoolConceptResponse
func ToSchoolConceptResponse(sc *entities.SchoolConcept) SchoolConceptResponse {
	return SchoolConceptResponse{
		ID:        sc.ID.String(),
		TermKey:   sc.TermKey,
		TermValue: sc.TermValue,
		Category:  sc.Category,
	}
}

// ToSchoolConceptResponseList converts a slice of SchoolConcept entities to SchoolConceptResponse slice
func ToSchoolConceptResponseList(concepts []*entities.SchoolConcept) []SchoolConceptResponse {
	responses := make([]SchoolConceptResponse, len(concepts))
	for i, sc := range concepts {
		responses[i] = ToSchoolConceptResponse(sc)
	}
	return responses
}
