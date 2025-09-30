package dto

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest represents the request structure for user registration
type RegisterRequest struct {
	Name  string `json:"name" binding:"required,min=10,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest represents the request structure for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=10,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

// UserResponse represents the response structure for user data
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersResponse represents the response structure for multiple users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}
