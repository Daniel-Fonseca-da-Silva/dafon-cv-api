package usecases

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase defines the interface for user business logic operations
type UserUseCase interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) (*dto.UsersResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// userUseCase implements UserUseCase interface
type userUseCase struct {
	userRepo          repositories.UserRepository
	configurationRepo repositories.ConfigurationRepository
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo repositories.UserRepository, configurationRepo repositories.ConfigurationRepository) UserUseCase {
	return &userUseCase{
		userRepo:          userRepo,
		configurationRepo: configurationRepo,
	}
}

// CreateUser creates a new user with business logic validation
func (uc *userUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user already exists with the same email
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save user to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create default configuration for the user
	configuration := &models.Configuration{
		UserID:        user.ID,
		Language:      "en-us", // Default language
		Newsletter:    false,   // Default: newsletter off
		ReceiveEmails: false,   // Default: receive emails off
	}

	if err := uc.configurationRepo.Create(ctx, configuration); err != nil {
		return nil, err
	}

	// Return response
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByID retrieves a user by ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetAllUsers retrieves all users
func (uc *userUseCase) GetAllUsers(ctx context.Context) (*dto.UsersResponse, error) {
	users, err := uc.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
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
		return nil, err
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
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
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

	return uc.userRepo.Delete(ctx, id)
}
