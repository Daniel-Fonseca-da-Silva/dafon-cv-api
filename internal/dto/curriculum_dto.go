package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateCurriculumRequest represents the request structure for creating a curriculum
type CreateCurriculumRequest struct {
	UserID        string                   `json:"user_id" binding:"required,uuid"`
	FullName      string                   `json:"full_name" binding:"required,min=5,max=50"`
	Email         string                   `json:"email" binding:"required,email"`
	Phone         string                   `json:"phone" binding:"required"`
	DriverLicense string                   `json:"driver_license" binding:"required"`
	Intro         string                   `json:"intro" binding:"required"`
	Skills        string                   `json:"skills" binding:"required"`
	Languages     string                   `json:"languages" binding:"required"`
	Courses       string                   `json:"courses"`
	SocialLinks   string                   `json:"social_links"`
	Works         []CreateWorkRequest      `json:"works"`
	Educations    []CreateEducationRequest `json:"educations"`
}

// CurriculumResponse represents the response structure for curriculum data
type CurriculumResponse struct {
	ID            uuid.UUID           `json:"id"`
	FullName      string              `json:"full_name"`
	Email         string              `json:"email"`
	Phone         string              `json:"phone"`
	DriverLicense string              `json:"driver_license"`
	Intro         string              `json:"intro"`
	Skills        string              `json:"skills"`
	Languages     string              `json:"languages"`
	Courses       string              `json:"courses"`
	SocialLinks   string              `json:"social_links"`
	Works         []WorkResponse      `json:"works"`
	Educations    []EducationResponse `json:"educations"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}
