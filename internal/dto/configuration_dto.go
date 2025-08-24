package dto

import (
	"time"

	"github.com/google/uuid"
)

// UpdateConfigurationRequest represents the request structure for updating a configuration
type UpdateConfigurationRequest struct {
	Language      string `json:"language" binding:"omitempty,min=2,max=255"`
	Newsletter    bool   `json:"newsletter" binding:"omitempty"`
	ReceiveEmails bool   `json:"receive_emails" binding:"omitempty"`
}

// ConfigurationResponse represents the response structure for configuration data
type ConfigurationResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Language      string    `json:"language"`
	Newsletter    bool      `json:"newsletter"`
	ReceiveEmails bool      `json:"receive_emails"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
