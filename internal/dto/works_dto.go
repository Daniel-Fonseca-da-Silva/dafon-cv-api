package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateWorkRequest represents the request structure for creating a work entry
type CreateWorkRequest struct {
	Position    string     `json:"position" binding:"required,min=2,max=255"`
	Company     string     `json:"company" binding:"required,min=2,max=255"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"` // Nullable para trabalhos atuais
}

// WorkResponse represents the response structure for work data
type WorkResponse struct {
	ID          uuid.UUID  `json:"id"`
	Position    string     `json:"position"`
	Company     string     `json:"company"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
