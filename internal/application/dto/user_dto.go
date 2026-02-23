package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts a User entity to UserResponse
func ToUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FirstName + " " + user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponseList converts a slice of User entities to UserResponse slice
func ToUserResponseList(users []*entities.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
