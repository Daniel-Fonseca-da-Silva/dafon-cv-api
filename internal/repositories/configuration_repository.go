package repositories

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
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
	db *gorm.DB
}

// NewConfigurationRepository creates a new instance of ConfigurationRepository
func NewConfigurationRepository(db *gorm.DB) ConfigurationRepository {
	return &configurationRepository{
		db: db,
	}
}

// Create creates a new configuration in the database
func (c *configurationRepository) Create(ctx context.Context, configuration *models.Configuration) error {
	return c.db.WithContext(ctx).Create(configuration).Error
}

// GetByID retrieves a configuration by ID
func (c *configurationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Configuration, error) {
	var configuration models.Configuration
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&configuration).Error
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

// GetByUserID retrieves a configuration by User ID
func (c *configurationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Configuration, error) {
	var configuration models.Configuration
	err := c.db.WithContext(ctx).Where("user_id = ?", userID).First(&configuration).Error
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

// Update updates an existing configuration in the database
func (c *configurationRepository) Update(ctx context.Context, configuration *models.Configuration) error {
	return c.db.WithContext(ctx).Save(configuration).Error
}

// Delete removes a configuration from the database
func (c *configurationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&models.Configuration{}, id).Error
}
