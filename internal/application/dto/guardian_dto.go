package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateGuardianRelationRequest represents the request to create a guardian relation
type CreateGuardianRelationRequest struct {
	GuardianID       string `json:"guardian_id" binding:"required"`
	StudentID        string `json:"student_id" binding:"required"`
	RelationshipType string `json:"relationship_type" binding:"required"`
}

// UpdateGuardianRelationRequest represents the request to update a guardian relation
type UpdateGuardianRelationRequest struct {
	RelationshipType *string `json:"relationship_type,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
}

// GuardianRelationResponse represents a guardian relation in API responses
type GuardianRelationResponse struct {
	ID               string    `json:"id"`
	GuardianID       string    `json:"guardian_id"`
	StudentID        string    `json:"student_id"`
	RelationshipType string    `json:"relationship_type"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedBy        string    `json:"created_by"`
}

// ToGuardianRelationResponse converts a GuardianRelation entity to response
func ToGuardianRelationResponse(relation *entities.GuardianRelation) *GuardianRelationResponse {
	return &GuardianRelationResponse{
		ID:               relation.ID.String(),
		GuardianID:       relation.GuardianID.String(),
		StudentID:        relation.StudentID.String(),
		RelationshipType: relation.RelationshipType,
		IsActive:         relation.IsActive,
		CreatedAt:        relation.CreatedAt,
		UpdatedAt:        relation.UpdatedAt,
		CreatedBy:        relation.CreatedBy,
	}
}

// ToGuardianRelationResponseList converts a slice of GuardianRelation entities to responses
func ToGuardianRelationResponseList(relations []*entities.GuardianRelation) []*GuardianRelationResponse {
	responses := make([]*GuardianRelationResponse, len(relations))
	for i, r := range relations {
		responses[i] = ToGuardianRelationResponse(r)
	}
	return responses
}
