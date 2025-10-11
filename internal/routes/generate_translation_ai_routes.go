package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/ratelimit"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGenerateTranslationAIRoutes configures AI filtering-related routes
func SetupGenerateTranslationAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateTranslationAIUseCase, err := usecases.NewGenerateTranslationAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Translation AI usecase", zap.Error(err))
		return
	}

	generateTranslationAIHandler := handlers.NewGenerateTranslationAIHandler(generateTranslationAIUseCase)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateTranslations := router.Group("/api/v1/generate-translation-ai")
	generateTranslations.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	generateTranslations.Use(ratelimit.RateLimiterMiddleware(aiRateLimiter))
	{
		generateTranslations.POST("", generateTranslationAIHandler.FilterContent)
	}
}
