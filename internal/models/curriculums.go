package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Curriculums struct {
	ID uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:curriculums"`
	gorm.Model
	FullName          string    `json:"full_name" gorm:"size:255;not null" validate:"required"`
	Email             string    `json:"email" gorm:"size:255;not null" validate:"required,email"`
	Phone             string    `json:"phone" gorm:"size:20;not null" validate:"required"`
	DriverLicense     string    `json:"driver_license" gorm:"size:255"`
	Intro             string    `json:"intro" gorm:"type:text;not null" validate:"required"`
	Technologies      string    `json:"technologies" gorm:"type:text;not null" validate:"required"`
	DateDisponibility time.Time `json:"date_disponibility" gorm:"type:date"`
	Languages         string    `json:"languages" gorm:"type:text;not null" validate:"required"`
	LevelEducation    string    `json:"level_education" gorm:"size:255;not null" validate:"required"`
	Courses           string    `json:"courses" gorm:"type:text"`
	SocialLinks       string    `json:"social_links" gorm:"size:100" validate:"url"`
	JobDescription    string    `json:"job_description" gorm:"type:text"`
	Works             []Work    `json:"works" gorm:"foreignKey:CurriculumID"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index" validate:"required"`
	User              User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Curriculums) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
