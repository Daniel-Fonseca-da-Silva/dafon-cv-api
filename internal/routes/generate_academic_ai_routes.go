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

// SetupGenerateAcademicAIRoutes configures AI filtering-related routes
func SetupGenerateAcademicAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateAcademicAIUseCase, err := usecases.NewGenerateAcademicAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Academic AI usecase", zap.Error(err))
		return
	}

	generateAcademicAIHandler := handlers.NewGenerateAcademicAIHandler(generateAcademicAIUseCase)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateAcademic := router.Group("/api/v1/generate-academic-ai")
	generateAcademic.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	generateAcademic.Use(ratelimit.RateLimiterMiddleware(aiRateLimiter))
	{
		generateAcademic.POST("", generateAcademicAIHandler.FilterContent)
	}
}
