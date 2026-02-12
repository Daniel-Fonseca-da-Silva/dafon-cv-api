package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID());table:sessions"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Token     string    `json:"token" gorm:"size:255;unique;not null;index"`
	IsActive  bool      `json:"is_active" gorm:"not null;default:true;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
