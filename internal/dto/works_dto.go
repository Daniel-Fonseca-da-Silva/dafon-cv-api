package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateWorkRequest represents the request structure for creating a work entry
type CreateWorkRequest struct {
	JobTitle           string     `json:"job_title" binding:"required,min=2,max=255"`
	CompanyName        string     `json:"company_name" binding:"required,min=2,max=255"`
	CompanyDescription string     `json:"company_description"`
	StartDate          time.Time  `json:"start_date" binding:"required"`
	EndDate            *time.Time `json:"end_date"` // Nullable para trabalhos atuais
}

// WorkResponse represents the response structure for work data
type WorkResponse struct {
	ID                 uuid.UUID  `json:"id"`
	JobTitle           string     `json:"job_title"`
	CompanyName        string     `json:"company_name"`
	CompanyDescription string     `json:"company_description"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            *time.Time `json:"end_date"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
