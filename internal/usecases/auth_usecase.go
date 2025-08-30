package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID string) error
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error)
}

// authUseCase implements AuthUseCase
type authUseCase struct {
	userRepo          repositories.UserRepository
	passwordResetRepo repositories.PasswordResetRepository
	userUseCase       UserUseCase
	emailUseCase      EmailUseCase
	jwtSecret         string
	jwtDuration       time.Duration
	appURL            string
}

// NewAuthUseCase creates a new instance of AuthUseCase
func NewAuthUseCase(userRepo repositories.UserRepository, passwordResetRepo repositories.PasswordResetRepository, userUseCase UserUseCase, emailUseCase EmailUseCase, jwtSecret string, jwtDuration time.Duration, appURL string) AuthUseCase {
	return &authUseCase{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		userUseCase:       userUseCase,
		emailUseCase:      emailUseCase,
		jwtSecret:         jwtSecret,
		jwtDuration:       jwtDuration,
		appURL:            appURL,
	}
}

// Login authenticates a user and returns a JWT token
func (uc *authUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
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
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
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

// ForgotPassword handles the forgot password request
func (uc *authUseCase) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	// Check if user exists
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if user exists or not for security reasons
		return &dto.ForgotPasswordResponse{
			Message: "Se o email existir em nossa base de dados, você receberá um link de recuperação de senha.",
		}, nil
	}

	// Delete any existing password reset tokens for this user
	uc.passwordResetRepo.DeleteByUserID(ctx, user.ID)

	// Generate secure token
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	// Create password reset record
	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		Used:      false,
	}

	if err := uc.passwordResetRepo.Create(ctx, passwordReset); err != nil {
		return nil, err
	}

	// Generate reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", uc.appURL, token)

	// Send email
	if err := uc.emailUseCase.SendPasswordResetEmail(user.Email, resetLink); err != nil {
		return nil, err
	}

	return &dto.ForgotPasswordResponse{
		Message: "Se o email existir em nossa base de dados, você receberá um link de recuperação de senha.",
	}, nil
}

// ResetPassword handles the password reset request
func (uc *authUseCase) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	// Get password reset by token
	passwordReset, err := uc.passwordResetRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, errors.ErrTokenExpired
	}

	// Check if token is already used
	if passwordReset.Used {
		return nil, errors.ErrInvalidToken
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, passwordReset.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update user password
	user.Password = string(hashedPassword)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Mark token as used
	if err := uc.passwordResetRepo.MarkAsUsed(ctx, req.Token); err != nil {
		return nil, err
	}

	// Delete all password reset tokens for this user
	uc.passwordResetRepo.DeleteByUserID(ctx, user.ID)

	return &dto.ResetPasswordResponse{
		Message: "Senha alterada com sucesso!",
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
