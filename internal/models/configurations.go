package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Configuration struct {
	gorm.Model
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:configuration"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	User       User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Language   string    `json:"language" gorm:"size:10;not null;default:'en-us'"`
	Newsletter bool      `json:"newsletter" gorm:"not null;default:false"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Configuration) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
