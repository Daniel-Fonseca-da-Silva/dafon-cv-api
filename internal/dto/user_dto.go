package dto

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest represents the request structure for user registration
type RegisterRequest struct {
	Name     string  `json:"name" binding:"required,min=15,max=100"`
	Email    string  `json:"email" binding:"required,email"`
	ImageURL *string `json:"image_url" binding:"omitempty,url"`
}

// UpdateUserRequest represents the request structure for updating a user
type UpdateUserRequest struct {
	Name       string   `json:"name" binding:"omitempty,min=10,max=100"`
	Email      string   `json:"email" binding:"omitempty,email"`
	ImageURL   *string  `json:"image_url" binding:"omitempty,url"`
	Country    string   `json:"country" binding:"omitempty,len=2"`
	State      string   `json:"state" binding:"omitempty,max=255"`
	City       string   `json:"city" binding:"omitempty,max=255"`
	Phone      string   `json:"phone" binding:"omitempty,phone"`
	Employment *bool    `json:"employment" binding:"omitempty"`
	Gender     string   `json:"gender" binding:"omitempty,oneof=male female"`
	Age        *int     `json:"age" binding:"omitempty,min=0"`
	Salary     *float64 `json:"salary" binding:"omitempty"`
	Migration  *bool    `json:"migration" binding:"omitempty"`
	Admin      *bool    `json:"admin" binding:"omitempty"`
}

// UserResponse represents the response structure for user data
type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ImageURL   *string   `json:"image_url,omitempty"`
	Country    string    `json:"country,omitempty"`
	State      string    `json:"state,omitempty"`
	City       string    `json:"city,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Employment bool      `json:"employment,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	Age        int       `json:"age,omitempty"`
	Salary     float64   `json:"salary,omitempty"`
	Migration  bool      `json:"migration,omitempty"`
	Admin      bool      `json:"admin,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UsersResponse represents the response structure for multiple users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

// SendEmailRequest represents the request structure for sending authentication email
type SendEmailRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	URLToken string `json:"url_token" binding:"required,min=1"`
}

// SendEmailResponse represents the response structure for sending authentication email
type SendEmailResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
