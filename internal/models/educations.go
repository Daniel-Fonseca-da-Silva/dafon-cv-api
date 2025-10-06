package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Education struct {
	gorm.Model
	ID           uuid.UUID  `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:educations"`
	CurriculumID uuid.UUID  `json:"curriculum_id" gorm:"type:char(36);not null;index"`
	Institution  string     `json:"institution" gorm:"size:255;not null"`
	Degree       string     `json:"degree" gorm:"size:255;not null"`
	StartDate    time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate      *time.Time `json:"end_date" gorm:"type:date"` // Nullable para trabalhos atuais
	Description  string     `json:"description" gorm:"type:text"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (e *Education) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
