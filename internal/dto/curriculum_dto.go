package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateCurriculumRequest represents the request structure for creating a curriculum
type CreateCurriculumRequest struct {
	FullName          string              `json:"full_name" binding:"required,min=5,max=50"`
	Email             string              `json:"email" binding:"required,email"`
	Phone             string              `json:"phone" binding:"required"`
	DriverLicense     string              `json:"driver_license" binding:"required"`
	Intro             string              `json:"intro" binding:"required"`
	Technologies      string              `json:"technologies" binding:"required"`
	DateDisponibility *time.Time          `json:"date_disponibility"`
	Languages         string              `json:"languages" binding:"required"`
	LevelEducation    string              `json:"level_education" binding:"required,min=5"`
	Courses           string              `json:"courses"`
	SocialLinks       string              `json:"social_links"`
	JobDescription    string              `json:"job_description"`
	Works             []CreateWorkRequest `json:"works"`
}

// CurriculumResponse represents the response structure for curriculum data
type CurriculumResponse struct {
	ID                uuid.UUID      `json:"id"`
	FullName          string         `json:"full_name"`
	Email             string         `json:"email"`
	Phone             string         `json:"phone"`
	DriverLicense     string         `json:"driver_license"`
	Intro             string         `json:"intro"`
	Technologies      string         `json:"technologies"`
	DateDisponibility time.Time      `json:"date_disponibility"`
	Languages         string         `json:"languages"`
	LevelEducation    string         `json:"level_education"`
	Courses           string         `json:"courses"`
	SocialLinks       string         `json:"social_links"`
	JobDescription    string         `json:"job_description"`
	Works             []WorkResponse `json:"works"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
