package dto

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateAcademicUnitRequest represents the request to create an academic unit
type CreateAcademicUnitRequest struct {
	ParentUnitID *string                `json:"parent_unit_id"`
	Type         string                 `json:"type" binding:"required"`
	DisplayName  string                 `json:"display_name" binding:"required,min=3"`
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// UpdateAcademicUnitRequest represents the request to update an academic unit
type UpdateAcademicUnitRequest struct {
	ParentUnitID *string                `json:"parent_unit_id"`
	DisplayName  *string                `json:"display_name"`
	Description  *string                `json:"description"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// AcademicUnitResponse represents an academic unit in API responses
type AcademicUnitResponse struct {
	ID           string                 `json:"id"`
	ParentUnitID *string                `json:"parent_unit_id,omitempty"`
	SchoolID     string                 `json:"school_id"`
	Type         string                 `json:"type"`
	DisplayName  string                 `json:"display_name"`
	Code         string                 `json:"code,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
}

// UnitTreeNode represents a node in the hierarchical tree
type UnitTreeNode struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	DisplayName string          `json:"display_name"`
	Code        string          `json:"code,omitempty"`
	Depth       int             `json:"depth"`
	Children    []*UnitTreeNode `json:"children,omitempty"`
}

// ToAcademicUnitResponse converts an AcademicUnit entity to AcademicUnitResponse
func ToAcademicUnitResponse(unit *entities.AcademicUnit) AcademicUnitResponse {
	var parentID *string
	if unit.ParentUnitID != nil {
		id := unit.ParentUnitID.String()
		parentID = &id
	}

	var metadata map[string]interface{}
	if len(unit.Metadata) > 0 {
		_ = json.Unmarshal(unit.Metadata, &metadata)
	}

	desc := ""
	if unit.Description != nil {
		desc = *unit.Description
	}

	return AcademicUnitResponse{
		ID:           unit.ID.String(),
		ParentUnitID: parentID,
		SchoolID:     unit.SchoolID.String(),
		Type:         unit.Type,
		DisplayName:  unit.Name,
		Code:         unit.Code,
		Description:  desc,
		Metadata:     metadata,
		CreatedAt:    unit.CreatedAt,
		UpdatedAt:    unit.UpdatedAt,
		DeletedAt:    unit.DeletedAt,
	}
}

// ToAcademicUnitResponseList converts a slice of AcademicUnit entities to responses
func ToAcademicUnitResponseList(units []*entities.AcademicUnit) []AcademicUnitResponse {
	responses := make([]AcademicUnitResponse, len(units))
	for i, unit := range units {
		responses[i] = ToAcademicUnitResponse(unit)
	}
	return responses
}

// BuildUnitTree builds a hierarchical tree from a flat list
func BuildUnitTree(units []*entities.AcademicUnit) []*UnitTreeNode {
	if len(units) == 0 {
		return []*UnitTreeNode{}
	}

	unitMap := make(map[string]*UnitTreeNode)
	var roots []*UnitTreeNode

	for _, unit := range units {
		node := &UnitTreeNode{
			ID:          unit.ID.String(),
			Type:        unit.Type,
			DisplayName: unit.Name,
			Code:        unit.Code,
			Depth:       1,
			Children:    []*UnitTreeNode{},
		}
		unitMap[unit.ID.String()] = node
	}

	for _, unit := range units {
		node := unitMap[unit.ID.String()]
		if unit.ParentUnitID == nil {
			roots = append(roots, node)
		} else {
			parentID := unit.ParentUnitID.String()
			if parent, exists := unitMap[parentID]; exists {
				node.Depth = parent.Depth + 1
				parent.Children = append(parent.Children, node)
			} else {
				roots = append(roots, node)
			}
		}
	}

	return roots
}
