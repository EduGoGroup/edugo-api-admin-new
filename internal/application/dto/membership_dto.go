package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateMembershipRequest represents the request to create a membership
type CreateMembershipRequest struct {
	UnitID string `json:"unit_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// UpdateMembershipRequest represents the request to update a membership
type UpdateMembershipRequest struct {
	Role *string `json:"role"`
}

// MembershipResponse represents a membership in API responses
type MembershipResponse struct {
	ID          string     `json:"id"`
	UnitID      string     `json:"unit_id"`
	UserID      string     `json:"user_id"`
	Role        string     `json:"role"`
	EnrolledAt  time.Time  `json:"enrolled_at"`
	WithdrawnAt *time.Time `json:"withdrawn_at,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToMembershipResponse converts a Membership entity to MembershipResponse
func ToMembershipResponse(m *entities.Membership) MembershipResponse {
	var unitID string
	if m.AcademicUnitID != nil {
		unitID = m.AcademicUnitID.String()
	}
	return MembershipResponse{
		ID:          m.ID.String(),
		UnitID:      unitID,
		UserID:      m.UserID.String(),
		Role:        m.Role,
		EnrolledAt:  m.EnrolledAt,
		WithdrawnAt: m.WithdrawnAt,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToMembershipResponseList converts a slice of Membership entities to responses
func ToMembershipResponseList(memberships []*entities.Membership) []MembershipResponse {
	responses := make([]MembershipResponse, len(memberships))
	for i, m := range memberships {
		responses[i] = ToMembershipResponse(m)
	}
	return responses
}
