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
func SetupGenerateIntroAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateIntroAIUseCase, err := usecases.NewGenerateIntroAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Intro AI usecase", zap.Error(err))
		return
	}

	generateIntroAIHandler := handlers.NewGenerateIntroAIHandler(generateIntroAIUseCase)

	generateIntros := router.Group("/api/v1/generate-intro-ai")
	generateIntros.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))

	{
		generateIntros.POST("", generateIntroAIHandler.FilterContent)
	}
}
