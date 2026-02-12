package repositories

import (
	"errors"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *models.Session) error
	GetByToken(token string) (*models.Session, error)
	GetByUserID(userID uuid.UUID) ([]*models.Session, error)
	GetActiveByUserID(userID uuid.UUID) ([]*models.Session, error)
	GetActiveByUserIDAndToken(userID uuid.UUID, token string) (*models.Session, error)
	Update(session *models.Session) error
	DeactivateByUserID(userID uuid.UUID) error
	DeactivateByToken(token string) error
	DeleteExpired() error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) GetByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetByUserID(userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	err := r.db.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionRepository) GetActiveByUserID(userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	err := r.db.Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionRepository) GetActiveByUserIDAndToken(userID uuid.UUID, token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("user_id = ? AND token = ? AND is_active = ? AND expires_at > ?", userID, token, true, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Update(session *models.Session) error {
	return r.db.Save(session).Error
}

func (r *sessionRepository) DeactivateByUserID(userID uuid.UUID) error {
	return r.db.Model(&models.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

func (r *sessionRepository) DeactivateByToken(token string) error {
	return r.db.Model(&models.Session{}).Where("token = ?", token).Update("is_active", false).Error
}

func (r *sessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}
