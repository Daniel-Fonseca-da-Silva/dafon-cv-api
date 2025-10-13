package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/cache"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserUseCase defines the interface for user business logic operations
type UserUseCase interface {
	CreateUser(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) (*dto.UsersResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// userUseCase implements UserUseCase interface
type userUseCase struct {
	userRepo          repositories.UserRepository
	configurationRepo repositories.ConfigurationRepository
	cacheService      *cache.CacheService
	logger            *zap.Logger
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo repositories.UserRepository, configurationRepo repositories.ConfigurationRepository, cacheService *cache.CacheService, logger *zap.Logger) UserUseCase {
	return &userUseCase{
		userRepo:          userRepo,
		configurationRepo: configurationRepo,
		cacheService:      cacheService,
		logger:            logger,
	}
}

// CreateUser creates a new user with business logic validation
func (uc *userUseCase) CreateUser(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists with the same email
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Create user model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		ImageURL: req.ImageURL,
	}

	// Save user to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to create user in database", zap.Error(err))
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	// Create default configuration for the user
	configuration := &models.Configuration{
		UserID:        user.ID,
		Language:      "en",  // Default language
		Newsletter:    false, // Default: newsletter off
		ReceiveEmails: false, // Default: receive emails off
	}

	if err := uc.configurationRepo.Create(ctx, configuration); err != nil {
		uc.logger.Error("Failed to create user configuration", zap.Error(err))
		return nil, fmt.Errorf("failed to create user configuration: %w", err)
	}

	// Return response
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByID retrieves a user by ID with cache support
func (uc *userUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	cacheKey := cache.GenerateUserCacheKey(id.String())

	// Try to get from cache first
	var userResponse dto.UserResponse
	found, err := uc.cacheService.Get(ctx, cacheKey, &userResponse)
	if err != nil {
		uc.logger.Warn("Failed to get user from cache, falling back to database",
			zap.Error(err),
			zap.String("user_id", id.String()))
	} else if found {
		uc.logger.Debug("User retrieved from cache")
		return &userResponse, nil
	}

	// Cache miss - get from database
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user by ID from database", zap.Error(err))
		return nil, fmt.Errorf("failed to get user by ID from database: %w", err)
	}

	// Create response
	userResponse = dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Armazena os dados em cache por 15 minutos
	ttl := 15 * time.Minute
	if err := uc.cacheService.Set(ctx, cacheKey, userResponse, ttl); err != nil {
		uc.logger.Warn("Failed to cache user data",
			zap.Error(err),
			zap.String("user_id", id.String()),
			zap.Duration("ttl", ttl))
	} else {
		uc.logger.Debug("User cached successfully",
			zap.String("user_id", id.String()),
			zap.Duration("ttl", ttl))
	}

	return &userResponse, nil
}

// GetAllUsers retrieves all users
func (uc *userUseCase) GetAllUsers(ctx context.Context) (*dto.UsersResponse, error) {
	users, err := uc.userRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all users from database", zap.Error(err))
		return nil, fmt.Errorf("failed to get all users from database: %w", err)
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			ImageURL:  user.ImageURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return &dto.UsersResponse{
		Users: userResponses,
		Total: len(userResponses),
	}, nil
}

// UpdateUser updates an existing user
func (uc *userUseCase) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user by ID for update", zap.Error(err), zap.String("user_id", id.String()))
		return nil, fmt.Errorf("failed to get user by ID for update: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already taken by another user
		if req.Email != user.Email {
			existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
			if err == nil && existingUser != nil {
				return nil, errors.ErrUserAlreadyExists
			}
		}
		user.Email = req.Email
	}
	if req.ImageURL != nil {
		user.ImageURL = req.ImageURL
	}

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user in database", zap.Error(err), zap.String("user_id", id.String()))
		return nil, fmt.Errorf("failed to update user in database: %w", err)
	}

	// Invalidate cache for this user
	cacheKey := cache.GenerateUserCacheKey(id.String())
	if err := uc.cacheService.Delete(ctx, cacheKey); err != nil {
		uc.logger.Warn("Failed to invalidate user cache after update",
			zap.Error(err),
			zap.String("user_id", id.String()))
	} else {
		uc.logger.Debug("User cache invalidated after update", zap.String("user_id", id.String()))
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// DeleteUser removes a user from the system
func (uc *userUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete user from database
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache for this user
	cacheKey := cache.GenerateUserCacheKey(id.String())
	if err := uc.cacheService.Delete(ctx, cacheKey); err != nil {
		uc.logger.Warn("Failed to invalidate user cache after deletion",
			zap.Error(err),
			zap.String("user_id", id.String()))
	} else {
		uc.logger.Debug("User cache invalidated after deletion", zap.String("user_id", id.String()))
	}

	return nil
}
