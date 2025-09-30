package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupEmailRoutes configures authentication email-related routes
func SetupEmailRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	// Initialize email dependencies
	emailUseCase, err := usecases.NewEmailUseCase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize email use case", zap.Error(err))
		return
	}

	emailHandler := handlers.NewEmailHandler(emailUseCase, logger)

	// Protected email routes (authentication required)
	email := router.Group("/api/v1/send-email")
	email.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	{
		email.POST("", emailHandler.SendEmail)
	}
}
