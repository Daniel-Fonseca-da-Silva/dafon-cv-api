package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) error {
	// Initialize JWT configuration
	jwtConfig, err := config.NewJWTConfig(logger)
	if err != nil {
		return err
	}

	// Initialize session configuration
	sessionConfig, err := config.NewSessionConfig(logger)
	if err != nil {
		return err
	}

	// Initialize auth dependencies
	userRepo := repositories.NewUserRepository(db)
	configurationRepo := repositories.NewConfigurationRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo)

	// Initialize session dependencies
	sessionRepo := repositories.NewSessionRepository(db)
	emailUseCase, err := usecases.NewEmailUseCase(logger)
	if err != nil {
		logger.Error("Failed to initialize email use case", zap.Error(err))
		return err
	}
	sessionUseCase := usecases.NewSessionUseCase(sessionRepo, emailUseCase, logger)

	// Initialize auth use case with all dependencies
	authUseCase := usecases.NewAuthUseCase(
		userRepo,
		userUseCase,
		sessionUseCase,
		emailUseCase,
		jwtConfig.SecretKey,
		jwtConfig.Duration,
		sessionConfig.Duration,
		cfg.App.URL,
		logger,
	)
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Auth routes group
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.GET("/login-with-token", authHandler.LoginWithToken)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", middleware.AuthMiddleware(jwtConfig), authHandler.Logout)
	}

	return nil
}
