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
	jwtConfig := config.NewJWTConfig(logger)

	// Initialize auth dependencies
	userRepo := repositories.NewUserRepository(db)
	configurationRepo := repositories.NewConfigurationRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo)

	authUseCase := usecases.NewAuthUseCase(userRepo, userUseCase, jwtConfig.SecretKey, jwtConfig.Duration)
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Auth routes group
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", middleware.AuthMiddleware(jwtConfig), authHandler.Logout)
	}

	return nil
}
