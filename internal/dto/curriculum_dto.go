package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateCurriculumRequest represents the request structure for creating a curriculum
type CreateCurriculumRequest struct {
	FullName          string              `json:"full_name" binding:"required,min=2,max=255"`
	Email             string              `json:"email" binding:"required,email"`
	DriverLicense     string              `json:"driver_license"`
	AboutMe           string              `json:"about_me"`
	DateDisponibility time.Time           `json:"date_disponibility" binding:"required"`
	Languages         string              `json:"languages" binding:"required"`
	LevelEducation    string              `json:"level_education" binding:"required,min=2,max=255"`
	CompanyInfo       string              `json:"company_info"`
	Works             []CreateWorkRequest `json:"works"`
}

// CurriculumResponse represents the response structure for curriculum data
type CurriculumResponse struct {
	ID                uuid.UUID      `json:"id"`
	FullName          string         `json:"full_name"`
	Email             string         `json:"email"`
	DriverLicense     string         `json:"driver_license"`
	AboutMe           string         `json:"about_me"`
	DateDisponibility time.Time      `json:"date_disponibility"`
	Languages         string         `json:"languages"`
	LevelEducation    string         `json:"level_education"`
	CompanyInfo       string         `json:"company_info"`
	Works             []WorkResponse `json:"works"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
