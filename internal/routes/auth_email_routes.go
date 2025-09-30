package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupAuthEmailRoutes configures authentication email-related routes
func SetupAuthEmailRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	// Initialize email dependencies
	emailUseCase, err := usecases.NewEmailUseCase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize email use case", zap.Error(err))
		return
	}

	authEmailHandler := handlers.NewAuthEmailHandler(emailUseCase, logger)

	// Public auth email routes (no authentication required)
	authEmail := router.Group("/api/v1/auth")
	{
		authEmail.POST("/send-email", authEmailHandler.SendAuthEmail)
	}
}
