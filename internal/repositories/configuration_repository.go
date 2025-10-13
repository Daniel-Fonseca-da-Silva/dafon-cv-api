package repositories

import (
	"context"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ConfigurationRepository defines the interface for configuration data operations
type ConfigurationRepository interface {
	Create(ctx context.Context, configuration *models.Configuration) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Configuration, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Configuration, error)
	Update(ctx context.Context, configuration *models.Configuration) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// configurationRepository implements ConfigurationRepository interface
type configurationRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewConfigurationRepository creates a new instance of ConfigurationRepository
func NewConfigurationRepository(db *gorm.DB, logger *zap.Logger) ConfigurationRepository {
	return &configurationRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new configuration in the database
func (c *configurationRepository) Create(ctx context.Context, configuration *models.Configuration) error {
	if err := c.db.WithContext(ctx).Create(configuration).Error; err != nil {
		c.logger.Error("Failed to create configuration", zap.Error(err), zap.String("user_id", configuration.UserID.String()))
		return fmt.Errorf("failed to create configuration: %w", err)
	}
	return nil
}

// GetByID retrieves a configuration by ID
func (c *configurationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Configuration, error) {
	var configuration models.Configuration
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&configuration).Error
	if err != nil {
		c.logger.Error("Failed to get configuration by ID", zap.Error(err), zap.String("config_id", id.String()))
		return nil, fmt.Errorf("failed to get configuration by ID %s: %w", id.String(), err)
	}
	return &configuration, nil
}

// GetByUserID retrieves a configuration by User ID
func (c *configurationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Configuration, error) {
	var configuration models.Configuration
	err := c.db.WithContext(ctx).Where("user_id = ?", userID).First(&configuration).Error
	if err != nil {
		c.logger.Error("Failed to get configuration by user ID", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get configuration by user ID %s: %w", userID.String(), err)
	}
	return &configuration, nil
}

// Update updates an existing configuration in the database
func (c *configurationRepository) Update(ctx context.Context, configuration *models.Configuration) error {
	if err := c.db.WithContext(ctx).Save(configuration).Error; err != nil {
		c.logger.Error("Failed to update configuration", zap.Error(err), zap.String("config_id", configuration.ID.String()))
		return fmt.Errorf("failed to update configuration: %w", err)
	}
	return nil
}

// Delete removes a configuration from the database
func (c *configurationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.db.WithContext(ctx).Delete(&models.Configuration{}, id).Error; err != nil {
		c.logger.Error("Failed to delete configuration", zap.Error(err), zap.String("config_id", id.String()))
		return fmt.Errorf("failed to delete configuration: %w", err)
	}
	return nil
}
