package dto

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateSchoolRequest represents the request to create a school
type CreateSchoolRequest struct {
	Name             string                 `json:"name" binding:"required,min=3"`
	Code             string                 `json:"code" binding:"required,min=3"`
	Address          string                 `json:"address"`
	City             string                 `json:"city"`
	Country          string                 `json:"country"`
	ContactEmail     string                 `json:"contact_email"`
	ContactPhone     string                 `json:"contact_phone"`
	SubscriptionTier string                 `json:"subscription_tier"`
	MaxTeachers      int                    `json:"max_teachers"`
	MaxStudents      int                    `json:"max_students"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// UpdateSchoolRequest represents the request to update a school
type UpdateSchoolRequest struct {
	Name             *string                `json:"name"`
	Address          *string                `json:"address"`
	City             *string                `json:"city"`
	Country          *string                `json:"country"`
	ContactEmail     *string                `json:"contact_email"`
	ContactPhone     *string                `json:"contact_phone"`
	SubscriptionTier *string                `json:"subscription_tier"`
	MaxTeachers      *int                   `json:"max_teachers"`
	MaxStudents      *int                   `json:"max_students"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// SchoolResponse represents a school in API responses
type SchoolResponse struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Code             string                 `json:"code"`
	Address          string                 `json:"address"`
	City             string                 `json:"city,omitempty"`
	Country          string                 `json:"country"`
	ContactEmail     string                 `json:"contact_email,omitempty"`
	ContactPhone     string                 `json:"contact_phone,omitempty"`
	SubscriptionTier string                 `json:"subscription_tier"`
	MaxTeachers      int                    `json:"max_teachers"`
	MaxStudents      int                    `json:"max_students"`
	IsActive         bool                   `json:"is_active"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// ToSchoolResponse converts a School entity to SchoolResponse
func ToSchoolResponse(school *entities.School) SchoolResponse {
	var email, phone, address, city string
	if school.Email != nil {
		email = *school.Email
	}
	if school.Phone != nil {
		phone = *school.Phone
	}
	if school.Address != nil {
		address = *school.Address
	}
	if school.City != nil {
		city = *school.City
	}

	var metadata map[string]interface{}
	if len(school.Metadata) > 0 {
		_ = json.Unmarshal(school.Metadata, &metadata)
	}

	return SchoolResponse{
		ID:               school.ID.String(),
		Name:             school.Name,
		Code:             school.Code,
		Address:          address,
		City:             city,
		Country:          school.Country,
		ContactEmail:     email,
		ContactPhone:     phone,
		SubscriptionTier: school.SubscriptionTier,
		MaxTeachers:      school.MaxTeachers,
		MaxStudents:      school.MaxStudents,
		IsActive:         school.IsActive,
		Metadata:         metadata,
		CreatedAt:        school.CreatedAt,
		UpdatedAt:        school.UpdatedAt,
	}
}

// ToSchoolResponseList converts a slice of School entities to SchoolResponse slice
func ToSchoolResponseList(schools []*entities.School) []SchoolResponse {
	responses := make([]SchoolResponse, len(schools))
	for i, school := range schools {
		responses[i] = ToSchoolResponse(school)
	}
	return responses
}
