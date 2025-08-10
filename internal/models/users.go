package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:user"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Email       string         `json:"email" gorm:"size:255;unique;not null"`
	Password    string         `json:"-" gorm:"size:255;not null"` // "-" hides password from JSON
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Curriculums []Curriculums  `json:"curriculums,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
