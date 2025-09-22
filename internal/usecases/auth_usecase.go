package usecases

import (
	"context"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	LoginWithToken(ctx context.Context, token string) (*dto.AuthResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID string) error
}

// authUseCase implements AuthUseCase
type authUseCase struct {
	userRepo        repositories.UserRepository
	userUseCase     UserUseCase
	sessionUC       SessionUseCase
	emailUC         EmailUseCase
	jwtSecret       string
	jwtDuration     time.Duration
	sessionDuration time.Duration
	appBaseURL      string
	logger          *zap.Logger
}

// NewAuthUseCase creates a new instance of AuthUseCase
func NewAuthUseCase(userRepo repositories.UserRepository, userUseCase UserUseCase, sessionUC SessionUseCase, emailUC EmailUseCase, jwtSecret string, jwtDuration time.Duration, sessionDuration time.Duration, appBaseURL string, logger *zap.Logger) AuthUseCase {
	return &authUseCase{
		userRepo:        userRepo,
		userUseCase:     userUseCase,
		sessionUC:       sessionUC,
		emailUC:         emailUC,
		jwtSecret:       jwtSecret,
		jwtDuration:     jwtDuration,
		sessionDuration: sessionDuration,
		appBaseURL:      appBaseURL,
		logger:          logger,
	}
}

// Login verifies email, creates session token and sends it by email
func (uc *authUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	uc.logger.Info("Login request received",
		zap.String("email", req.Email),
	)

	// Check if user exists
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		uc.logger.Warn("Login attempt with non-existent email",
			zap.String("email", req.Email),
		)
		return nil, errors.ErrInvalidCredentials
	}

	// Create session token
	session, err := uc.sessionUC.CreateSession(user.ID, uc.sessionDuration)
	if err != nil {
		uc.logger.Error("Failed to create session for login",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	// Send session token by email
	if err := uc.emailUC.SendSessionTokenEmail(user.Email, session.Token, uc.appBaseURL); err != nil {
		uc.logger.Error("Failed to send session token email",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("Session token sent by email successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
		zap.String("session_id", session.ID.String()),
	)

	return &dto.LoginResponse{
		Message: "Login link sent to your email. Check your inbox and click the link to access your account.",
	}, nil
}

// Register creates a new user and sends a session token link
func (uc *authUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	uc.logger.Info("Register request received",
		zap.String("email", req.Email),
		zap.String("name", req.Name),
	)

	// Create user using the user usecase (which will also create default configuration)
	user, err := uc.userUseCase.CreateUser(ctx, req)
	if err != nil {
		uc.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}

	// Create session token for new user
	session, err := uc.sessionUC.CreateSession(user.ID, uc.sessionDuration)
	if err != nil {
		uc.logger.Error("Failed to create session for new user",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	// Send session token email
	if err := uc.emailUC.SendSessionTokenEmail(user.Email, session.Token, uc.appBaseURL); err != nil {
		uc.logger.Error("Failed to send session token email to new user",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("User registered and session token sent successfully",
		zap.String("email", req.Email),
		zap.String("user_id", user.ID.String()),
		zap.String("session_id", session.ID.String()),
	)

	return &dto.AuthResponse{
		User: *user,
	}, nil
}

// Logout handles user logout by deactivating all sessions
func (uc *authUseCase) Logout(ctx context.Context, userID string) error {
	uc.logger.Info("Logout request received",
		zap.String("user_id", userID),
	)

	// Parse the user ID to validate it
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		uc.logger.Error("Invalid user ID format",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return errors.ErrInvalidToken
	}

	// Check if user exists
	_, err = uc.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		uc.logger.Error("User not found for logout",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return errors.ErrUserNotFound
	}

	// Deactivate all sessions for the user
	if err := uc.sessionUC.DeactivateUserSessions(userUUID); err != nil {
		uc.logger.Error("Failed to deactivate user sessions",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	uc.logger.Info("User logged out successfully",
		zap.String("user_id", userID),
	)

	return nil
}

// LoginWithToken validates session token and returns JWT
func (uc *authUseCase) LoginWithToken(ctx context.Context, token string) (*dto.AuthResponse, error) {
	uc.logger.Info("Login with token request received",
		zap.String("token", token),
	)

	// Validate the session token
	session, err := uc.sessionUC.ValidateSession(token)
	if err != nil {
		uc.logger.Error("Invalid or expired session token",
			zap.String("token", token),
			zap.Error(err),
		)
		return nil, err
	}

	// Get user information
	user, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		uc.logger.Error("User not found for token validation",
			zap.String("user_id", session.UserID.String()),
			zap.String("token", token),
			zap.Error(err),
		)
		return nil, errors.ErrUserNotFound
	}

	// Generate JWT token for the session
	jwtToken, expiresAt, err := uc.generateJWT(user.ID.String())
	if err != nil {
		uc.logger.Error("Failed to generate JWT for login",
			zap.String("user_id", user.ID.String()),
			zap.String("token", token),
			zap.Error(err),
		)
		return nil, err
	}

	// Deactivate the session token (one-time use)
	if err := uc.sessionUC.DeactivateSession(token); err != nil {
		uc.logger.Error("Failed to deactivate session token",
			zap.String("user_id", user.ID.String()),
			zap.String("token", token),
			zap.Error(err),
		)
		// Don't return error here, as the user is already authenticated
	}

	uc.logger.Info("Login completed successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return &dto.AuthResponse{
		Token:     &jwtToken,
		ExpiresAt: &expiresAt,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
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
