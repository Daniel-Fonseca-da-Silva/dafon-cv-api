package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the request structure for creating a user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest represents the request structure for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=255"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
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
