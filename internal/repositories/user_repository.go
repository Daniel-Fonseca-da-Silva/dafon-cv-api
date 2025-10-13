package repositories

import (
	"context"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("Failed to create user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		r.logger.Error("Failed to get user by ID", zap.Error(err), zap.String("user_id", id.String()))
		return nil, fmt.Errorf("failed to get user by ID %s: %w", id.String(), err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		r.logger.Error("Failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}
	return &user, nil
}

// GetAll retrieves all users from the database
func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		r.logger.Error("Failed to get all users", zap.Error(err))
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

// Update updates an existing user in the database
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error("Failed to update user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete removes a user from the database
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		r.logger.Error("Failed to delete user", zap.Error(err), zap.String("user_id", id.String()))
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
