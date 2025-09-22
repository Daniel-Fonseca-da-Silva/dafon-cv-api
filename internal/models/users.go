package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:user"`
	Name          string         `json:"name" gorm:"size:255;not null"`
	Email         string         `json:"email" gorm:"size:255;unique;not null"`
	Curriculums   []Curriculums  `json:"curriculums,omitempty" gorm:"foreignKey:UserID"`
	Configuration *Configuration `json:"configuration,omitempty" gorm:"foreignKey:UserID"`
	Sessions      []Session      `json:"sessions,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
