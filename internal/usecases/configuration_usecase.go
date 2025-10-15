package usecases

import (
	"context"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/cache"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConfigurationUseCase interface {
	GetConfigurationByUserID(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error)
	CreateDefaultConfiguration(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error)
	UpdateConfiguration(ctx context.Context, id uuid.UUID, req *dto.UpdateConfigurationRequest) (*dto.ConfigurationResponse, error)
	DeleteConfiguration(ctx context.Context, id uuid.UUID) error
}

type configurationUseCase struct {
	configurationRepo repositories.ConfigurationRepository
	cacheService      *cache.CacheService
	logger            *zap.Logger
}

func NewConfigurationUseCase(configurationRepo repositories.ConfigurationRepository, cacheService *cache.CacheService, logger *zap.Logger) ConfigurationUseCase {
	return &configurationUseCase{
		configurationRepo: configurationRepo,
		cacheService:      cacheService,
		logger:            logger,
	}
}

// GetConfigurationByUserID retrieves a configuration by user ID with cache support
func (c *configurationUseCase) GetConfigurationByUserID(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error) {
	cacheKey := cache.GenerateConfigurationCacheKey(userID.String())

	// Try to get from cache first
	var configurationResponse dto.ConfigurationResponse
	found, err := c.cacheService.Get(ctx, cacheKey, &configurationResponse)
	if err != nil {
		c.logger.Warn("Failed to get configuration from cache, falling back to database",
			zap.Error(err),
			zap.String("user_id", userID.String()))
	} else if found {
		c.logger.Debug("Configuration retrieved from cache", zap.String("user_id", userID.String()))
		return &configurationResponse, nil
	}

	// Cache miss - get from database
	configuration, err := c.configurationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create response
	configurationResponse = dto.ConfigurationResponse{
		ID:         configuration.ID,
		UserID:     configuration.UserID,
		Language:   configuration.Language,
		Newsletter: configuration.Newsletter,
		CreatedAt:  configuration.CreatedAt,
		UpdatedAt:  configuration.UpdatedAt,
	}

	// Armazena os dados em cache por 15 minutos
	ttl := 15 * time.Minute
	if err := c.cacheService.Set(ctx, cacheKey, configurationResponse, ttl); err != nil {
		c.logger.Warn("Failed to cache configuration data",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.Duration("ttl", ttl))
	} else {
		c.logger.Debug("Configuration cached successfully",
			zap.String("user_id", userID.String()),
			zap.Duration("ttl", ttl))
	}

	return &configurationResponse, nil
}

// CreateDefaultConfiguration creates a default configuration for a user
func (c *configurationUseCase) CreateDefaultConfiguration(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error) {
	configuration := &models.Configuration{
		UserID:     userID,
		Language:   "en-us", // Default language
		Newsletter: false,   // Default: newsletter off
	}

	if err := c.configurationRepo.Create(ctx, configuration); err != nil {
		return nil, err
	}

	return &dto.ConfigurationResponse{
		ID:         configuration.ID,
		UserID:     configuration.UserID,
		Language:   configuration.Language,
		Newsletter: configuration.Newsletter,
		CreatedAt:  configuration.CreatedAt,
		UpdatedAt:  configuration.UpdatedAt,
	}, nil
}

// UpdateConfiguration updates an existing configuration
func (c *configurationUseCase) UpdateConfiguration(ctx context.Context, id uuid.UUID, req *dto.UpdateConfigurationRequest) (*dto.ConfigurationResponse, error) {
	// Get existing configuration
	configuration, err := c.configurationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Language != "" {
		configuration.Language = req.Language
	}

	if req.Newsletter != configuration.Newsletter {
		configuration.Newsletter = req.Newsletter
	}

	// Save updated configuration
	if err := c.configurationRepo.Update(ctx, configuration); err != nil {
		return nil, err
	}

	// Invalidate cache for this user's configuration
	cacheKey := cache.GenerateConfigurationCacheKey(configuration.UserID.String())
	if err := c.cacheService.Delete(ctx, cacheKey); err != nil {
		c.logger.Warn("Failed to invalidate configuration cache after update",
			zap.Error(err),
			zap.String("user_id", configuration.UserID.String()))
	} else {
		c.logger.Debug("Configuration cache invalidated after update", zap.String("user_id", configuration.UserID.String()))
	}

	return &dto.ConfigurationResponse{
		ID:         configuration.ID,
		UserID:     configuration.UserID,
		Language:   configuration.Language,
		Newsletter: configuration.Newsletter,
		CreatedAt:  configuration.CreatedAt,
		UpdatedAt:  configuration.UpdatedAt,
	}, nil
}

// DeleteConfiguration deletes a configuration
func (c *configurationUseCase) DeleteConfiguration(ctx context.Context, id uuid.UUID) error {
	// Check if configuration exists and get user ID for cache invalidation
	configuration, err := c.configurationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete configuration from database
	if err := c.configurationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache for this user's configuration
	cacheKey := cache.GenerateConfigurationCacheKey(configuration.UserID.String())
	if err := c.cacheService.Delete(ctx, cacheKey); err != nil {
		c.logger.Warn("Failed to invalidate configuration cache after deletion",
			zap.Error(err),
			zap.String("user_id", configuration.UserID.String()))
	} else {
		c.logger.Debug("Configuration cache invalidated after deletion", zap.String("user_id", configuration.UserID.String()))
	}

	return nil
}
