package usecases

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
)

type ConfigurationUseCase interface {
	GetConfigurationByUserID(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error)
	CreateDefaultConfiguration(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error)
	UpdateConfiguration(ctx context.Context, id uuid.UUID, req *dto.UpdateConfigurationRequest) (*dto.ConfigurationResponse, error)
	DeleteConfiguration(ctx context.Context, id uuid.UUID) error
}

type configurationUseCase struct {
	configurationRepo repositories.ConfigurationRepository
}

func NewConfigurationUseCase(configurationRepo repositories.ConfigurationRepository) ConfigurationUseCase {
	return &configurationUseCase{
		configurationRepo: configurationRepo,
	}
}

// GetConfigurationByUserID retrieves a configuration by user ID
func (c *configurationUseCase) GetConfigurationByUserID(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error) {
	configuration, err := c.configurationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ConfigurationResponse{
		ID:            configuration.ID,
		UserID:        configuration.UserID,
		Language:      configuration.Language,
		Newsletter:    configuration.Newsletter,
		ReceiveEmails: configuration.ReceiveEmails,
		CreatedAt:     configuration.CreatedAt,
		UpdatedAt:     configuration.UpdatedAt,
	}, nil
}

// CreateDefaultConfiguration creates a default configuration for a user
func (c *configurationUseCase) CreateDefaultConfiguration(ctx context.Context, userID uuid.UUID) (*dto.ConfigurationResponse, error) {
	configuration := &models.Configuration{
		UserID:        userID,
		Language:      "en-us", // Default language
		Newsletter:    false,   // Default: newsletter off
		ReceiveEmails: false,   // Default: receive emails off
	}

	if err := c.configurationRepo.Create(ctx, configuration); err != nil {
		return nil, err
	}

	return &dto.ConfigurationResponse{
		ID:            configuration.ID,
		UserID:        configuration.UserID,
		Language:      configuration.Language,
		Newsletter:    configuration.Newsletter,
		ReceiveEmails: configuration.ReceiveEmails,
		CreatedAt:     configuration.CreatedAt,
		UpdatedAt:     configuration.UpdatedAt,
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

	if req.ReceiveEmails != configuration.ReceiveEmails {
		configuration.ReceiveEmails = req.ReceiveEmails
	}

	// Save updated configuration
	if err := c.configurationRepo.Update(ctx, configuration); err != nil {
		return nil, err
	}

	return &dto.ConfigurationResponse{
		ID:            configuration.ID,
		UserID:        configuration.UserID,
		Language:      configuration.Language,
		Newsletter:    configuration.Newsletter,
		ReceiveEmails: configuration.ReceiveEmails,
		CreatedAt:     configuration.CreatedAt,
		UpdatedAt:     configuration.UpdatedAt,
	}, nil
}

// DeleteConfiguration deletes a configuration
func (c *configurationUseCase) DeleteConfiguration(ctx context.Context, id uuid.UUID) error {
	// Check if configuration exists
	_, err := c.configurationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return c.configurationRepo.Delete(ctx, id)
}
