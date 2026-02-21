package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateSubjectRequest represents the request to create a subject
type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description"`
	Metadata    string `json:"metadata"`
}

// UpdateSubjectRequest represents the request to update a subject
type UpdateSubjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Metadata    *string `json:"metadata"`
}

// SubjectResponse represents a subject in API responses
type SubjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Metadata    string    `json:"metadata,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToSubjectResponse converts a Subject entity to SubjectResponse
func ToSubjectResponse(subject *entities.Subject) SubjectResponse {
	desc := ""
	if subject.Description != nil {
		desc = *subject.Description
	}
	meta := ""
	if subject.Metadata != nil {
		meta = *subject.Metadata
	}
	return SubjectResponse{
		ID:          subject.ID.String(),
		Name:        subject.Name,
		Description: desc,
		Metadata:    meta,
		IsActive:    subject.IsActive,
		CreatedAt:   subject.CreatedAt,
		UpdatedAt:   subject.UpdatedAt,
	}
}

// ToSubjectResponseList converts a slice of Subject entities to responses
func ToSubjectResponseList(subjects []*entities.Subject) []SubjectResponse {
	responses := make([]SubjectResponse, len(subjects))
	for i, s := range subjects {
		responses[i] = ToSubjectResponse(s)
	}
	return responses
}
