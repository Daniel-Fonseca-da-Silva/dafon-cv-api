package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Curriculums struct {
	ID                uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:curriculums"`
	FullName          string         `json:"full_name" gorm:"size:255;not null"`
	Email             string         `json:"email" gorm:"size:255;not null"`
	DriverLicense     string         `json:"driver_license" gorm:"size:255"`
	AboutMe           string         `json:"about_me" gorm:"type:text"`
	DateDisponibility time.Time      `json:"date_disponibility" gorm:"type:date"`
	Languages         string         `json:"languages" gorm:"type:text;not null"`
	LevelEducation    string         `json:"level_education" gorm:"size:255;not null"`
	CompanyInfo       string         `json:"company_info" gorm:"type:text"`
	Works             []Work         `json:"works" gorm:"foreignKey:CurriculumID"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Curriculums) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
