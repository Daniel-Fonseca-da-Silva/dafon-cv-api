package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Work struct {
	gorm.Model
	ID           uuid.UUID  `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:works"`
	CurriculumID uuid.UUID  `json:"curriculum_id" gorm:"type:char(36);not null;index"`
	Position     string     `json:"position" gorm:"size:255;not null"`
	Company      string     `json:"company" gorm:"size:255;not null"`
	Description  string     `json:"description" gorm:"type:text"`
	StartDate    time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate      *time.Time `json:"end_date" gorm:"type:date"` // Nullable para trabalhos atuais
}

// BeforeCreate will set a UUID rather than numeric ID
func (w *Work) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
