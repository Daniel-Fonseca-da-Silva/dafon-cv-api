package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGenerateTaskAIRoutes configures AI filtering-related routes
func SetupGenerateTaskAIRoutes(router *gin.Engine, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig := config.NewJWTConfig(logger)

	generateTaskAIUseCase, err := usecases.NewGenerateTaskAIUseCase()
	if err != nil {
		logger.Error("Failed to create GenerateTaskAI usecase", zap.Error(err))
		return
	}

	generateTaskAIHandler := handlers.NewGenerateTaskAIHandler(generateTaskAIUseCase)

	generateTasks := router.Group("/api/v1/generate-task-ai")
	generateTasks.Use(middleware.AuthMiddleware(jwtConfig))
	{
		generateTasks.POST("", generateTaskAIHandler.FilterContent)
	}
}
