package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGenerateIntroAIRoutes configures AI filtering-related routes
func SetupGenerateAnalyzeAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateAnalyzeAIUseCase, err := usecases.NewGenerateAnalyzeAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Analyze AI usecase", zap.Error(err))
		return
	}

	generateAnalyzeAIHandler := handlers.NewGenerateAnalyzeAIHandler(generateAnalyzeAIUseCase)

	generateAnalyze := router.Group("/api/v1/generate-analyze-ai")
	generateAnalyze.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))

	{
		generateAnalyze.POST("", generateAnalyzeAIHandler.FilterContent)
	}
}
