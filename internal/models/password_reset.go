package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordReset represents a password reset token
type PasswordReset struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	Token     string    `json:"token" gorm:"size:255;unique;not null"`
	Email     string    `json:"email" gorm:"size:255;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}
