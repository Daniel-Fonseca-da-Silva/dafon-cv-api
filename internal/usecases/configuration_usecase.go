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
	UpdateConfiguration(ctx context.Context, userID uuid.UUID, req *dto.UpdateConfigurationRequest) (*dto.ConfigurationResponse, error)
	DeleteConfiguration(ctx context.Context, userID uuid.UUID) error
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

// Retorna uma configuração por ID do usuário com suporte a cache
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

// Cria uma configuração padrão para um usuário
func (c *configurationUseCase) CreateDefaultConfiguration(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error) {
	configuration := &models.Configuration{
		UserID:     userID,
		Language:   "en-us", // Idioma padrão
		Newsletter: false,   // Newsletter padrão: off
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

// Atualiza uma configuração existente
func (c *configurationUseCase) UpdateConfiguration(ctx context.Context, userID uuid.UUID, req *dto.UpdateConfigurationRequest) (*dto.ConfigurationResponse, error) {
	// Obtém a configuração existente pelo ID do usuário
	configuration, err := c.configurationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Atualiza os campos se fornecidos
	if req.Language != "" {
		configuration.Language = req.Language
	}

	if req.Newsletter != configuration.Newsletter {
		configuration.Newsletter = req.Newsletter
	}

	// Salva a configuração atualizada
	if err := c.configurationRepo.Update(ctx, configuration); err != nil {
		return nil, err
	}

	// Invalida o cache para a configuração deste usuário
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

// Deleta uma configuração
func (c *configurationUseCase) DeleteConfiguration(ctx context.Context, userID uuid.UUID) error {
	// Verifica se a configuração existe e obtém a configuração para invalidar o cache
	configuration, err := c.configurationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Deleta a configuração do banco de dados
	if err := c.configurationRepo.Delete(ctx, configuration.ID); err != nil {
		return err
	}

	// Invalida o cache para a configuração deste usuário
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
