package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
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
	jwtConfig := config.NewJWTConfig(logger)

	// Initialize auth dependencies
	userRepo := repositories.NewUserRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	configurationRepo := repositories.NewConfigurationRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo)

	emailUseCase, err := usecases.NewEmailUseCase(logger)
	if err != nil {
		return errors.WrapError(err, "failed to initialize email use case")
	}

	authUseCase := usecases.NewAuthUseCase(userRepo, passwordResetRepo, userUseCase, emailUseCase, jwtConfig.SecretKey, jwtConfig.Duration, cfg.App.URL)
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Auth routes group
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", middleware.AuthMiddleware(jwtConfig), authHandler.Logout)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	return nil
}
