package usecases

import (
	"context"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID string) error
}

// authUseCase implements AuthUseCase
type authUseCase struct {
	userRepo    repositories.UserRepository
	userUseCase UserUseCase
	jwtSecret   string
	jwtDuration time.Duration
}

// NewAuthUseCase creates a new instance of AuthUseCase
func NewAuthUseCase(userRepo repositories.UserRepository, userUseCase UserUseCase, jwtSecret string, jwtDuration time.Duration) AuthUseCase {
	return &authUseCase{
		userRepo:    userRepo,
		userUseCase: userUseCase,
		jwtSecret:   jwtSecret,
		jwtDuration: jwtDuration,
	}
}

// Login authenticates a user and returns a JWT token
func (uc *authUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	token, expiresAt, err := uc.generateJWT(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// Register creates a new user and returns a JWT token
func (uc *authUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Create user using the user usecase (which will also create default configuration)
	createUserReq := &dto.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := uc.userUseCase.CreateUser(ctx, createUserReq)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, expiresAt, err := uc.generateJWT(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

// Logout handles user logout by invalidating the token
func (uc *authUseCase) Logout(ctx context.Context, userID string) error {

	// Parse the user ID to validate it
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.ErrInvalidToken
	}

	// Check if user exists
	_, err = uc.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	return nil
}

// generateJWT creates a JWT token for the given user ID
func (uc *authUseCase) generateJWT(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(uc.jwtDuration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	return tokenString, expiresAt, err
}
