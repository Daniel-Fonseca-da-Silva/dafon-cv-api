package dto

import "time"

// RegisterRequest represents the request structure for user registration
type RegisterRequest struct {
	Name  string `json:"name" binding:"required,min=10,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// AuthResponse represents the response structure for authentication
type AuthResponse struct {
	Token     *string      `json:"token,omitempty"`
	ExpiresAt *time.Time   `json:"expires_at,omitempty"`
	User      UserResponse `json:"user"`
}

// LogoutResponse represents the response structure for logout
type LogoutResponse struct {
	Message string `json:"message"`
}

// LoginResponse represents the response structure for login
type LoginResponse struct {
	Message string `json:"message"`
}

// LoginRequest represents the request structure for user login
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}
