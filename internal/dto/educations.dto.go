package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateEducationRequest struct {
	Institution string     `json:"institution" binding:"required,min=2,max=255"`
	Degree      string     `json:"degree" binding:"required,min=2,max=255"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	Description string     `json:"description" binding:"required,min=2,max=255"`
}

type EducationResponse struct {
	ID          uuid.UUID  `json:"id"`
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
