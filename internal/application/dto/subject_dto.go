package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateSubjectRequest represents the request to create a subject
type CreateSubjectRequest struct {
	Name           string `json:"name" binding:"required,min=2"`
	Description    string `json:"description"`
	AcademicUnitID string `json:"academic_unit_id"`
	Code           string `json:"code"`
}

// UpdateSubjectRequest represents the request to update a subject
type UpdateSubjectRequest struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	AcademicUnitID *string `json:"academic_unit_id"`
	Code           *string `json:"code"`
}

// SubjectResponse represents a subject in API responses
type SubjectResponse struct {
	ID             string  `json:"id"`
	SchoolID       string  `json:"school_id"`
	AcademicUnitID *string `json:"academic_unit_id,omitempty"`
	Name           string  `json:"name"`
	Code           *string `json:"code,omitempty"`
	Description    string  `json:"description,omitempty"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ToSubjectResponse converts a Subject entity to SubjectResponse
func ToSubjectResponse(subject *entities.Subject) SubjectResponse {
	desc := ""
	if subject.Description != nil {
		desc = *subject.Description
	}
	var academicUnitID *string
	if subject.AcademicUnitID != nil {
		s := subject.AcademicUnitID.String()
		academicUnitID = &s
	}
	return SubjectResponse{
		ID:             subject.ID.String(),
		SchoolID:       subject.SchoolID.String(),
		AcademicUnitID: academicUnitID,
		Name:           subject.Name,
		Code:           subject.Code,
		Description:    desc,
		IsActive:       subject.IsActive,
		CreatedAt:      subject.CreatedAt,
		UpdatedAt:      subject.UpdatedAt,
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
