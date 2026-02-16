package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CurriculumCreationStats holds the total number of successful curriculum creations per user.
// One row per user; TotalCreations only increments and is not affected by curriculum deletions.
type CurriculumCreationStats struct {
	gorm.Model
	UserID         uuid.UUID `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	TotalCreations int64     `json:"total_creations" gorm:"not null;default:0"`
	User           User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}
