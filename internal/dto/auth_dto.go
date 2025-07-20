package dto

import "time"

// LoginRequest represents the request structure for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the request structure for user registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse represents the response structure for authentication
type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// LogoutResponse represents the response structure for logout
type LogoutResponse struct {
	Message string `json:"message"`
}
