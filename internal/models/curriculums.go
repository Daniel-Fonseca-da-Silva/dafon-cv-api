package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Curriculums struct {
	ID uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:curriculums"`
	gorm.Model
	FullName      string      `json:"full_name" gorm:"size:255;not null" validate:"required"`
	Email         string      `json:"email" gorm:"size:255;not null" validate:"required,email"`
	Phone         string      `json:"phone" gorm:"size:20;not null" validate:"required"`
	DriverLicense string      `json:"driver_license" gorm:"size:255"`
	Intro         string      `json:"intro" gorm:"type:text;not null" validate:"required"`
	Skills        string      `json:"skills" gorm:"type:text;not null" validate:"required"`
	Languages     string      `json:"languages" gorm:"type:text;not null" validate:"required"`
	Courses       string      `json:"courses" gorm:"type:text"`
	SocialLinks   string      `json:"social_links" gorm:"type:text"`
	Works         []Work      `json:"works" gorm:"foreignKey:CurriculumID"`
	Educations    []Education `json:"educations" gorm:"foreignKey:CurriculumID"`
	UserID        uuid.UUID   `json:"user_id" gorm:"type:char(36);not null;index" validate:"required"`
	User          User        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Curriculums) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
