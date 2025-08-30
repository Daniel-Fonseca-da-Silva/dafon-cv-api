package repositories

import (
	"context"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordResetRepository defines the interface for password reset data operations
type PasswordResetRepository interface {
	Create(ctx context.Context, passwordReset *models.PasswordReset) error
	GetByToken(ctx context.Context, token string) (*models.PasswordReset, error)
	GetByEmail(ctx context.Context, email string) (*models.PasswordReset, error)
	MarkAsUsed(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// passwordResetRepository implements PasswordResetRepository interface
type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new instance of PasswordResetRepository
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{
		db: db,
	}
}

// Create creates a new password reset token in the database
func (r *passwordResetRepository) Create(ctx context.Context, passwordReset *models.PasswordReset) error {
	return r.db.WithContext(ctx).Create(passwordReset).Error
}

// GetByToken retrieves a password reset by token
func (r *passwordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ? AND used = ?", token, time.Now(), false).First(&passwordReset).Error
	if err != nil {
		return nil, err
	}
	return &passwordReset, nil
}

// GetByEmail retrieves a password reset by email
func (r *passwordResetRepository) GetByEmail(ctx context.Context, email string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	err := r.db.WithContext(ctx).Where("email = ? AND expires_at > ? AND used = ?", email, time.Now(), false).First(&passwordReset).Error
	if err != nil {
		return nil, err
	}
	return &passwordReset, nil
}

// MarkAsUsed marks a password reset token as used
func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.PasswordReset{}).Where("token = ?", token).Update("used", true).Error
}

// DeleteExpired deletes expired password reset tokens
func (r *passwordResetRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{}).Error
}

// DeleteByUserID deletes all password reset tokens for a specific user
func (r *passwordResetRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.PasswordReset{}).Error
}
