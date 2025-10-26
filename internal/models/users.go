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
	ImageURL      *string        `json:"image_url,omitempty" gorm:"size:512"`
	Country       string         `json:"country,omitempty" gorm:"size:2"`
	State         string         `json:"state,omitempty" gorm:"size:255"`
	City          string         `json:"city,omitempty" gorm:"size:255"`
	Phone         string         `json:"phone,omitempty" gorm:"size:20"`
	Employment    bool           `json:"employment,omitempty" gorm:"default:false"`
	Gender        string         `json:"gender,omitempty" gorm:"size:6"`
	Age           int            `json:"age,omitempty" gorm:"default:0"`
	Salary        float64        `json:"salary,omitempty" gorm:"type:double;default:0"`
	Migration     bool           `json:"migration,omitempty" gorm:"not null;default:false"`
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
