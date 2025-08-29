package dto

import "time"

// LoginRequest represents the request structure for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the request structure for user registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=10,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

// ForgotPasswordRequest represents the request structure for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request structure for reset password
type ResetPasswordRequest struct {
	Token    string `json:"token"` // Token is received via query parameter, not body
	Password string `json:"password" binding:"required,min=8,max=50"`
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

// ForgotPasswordResponse represents the response structure for forgot password
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordResponse represents the response structure for reset password
type ResetPasswordResponse struct {
	Message string `json:"message"`
}
