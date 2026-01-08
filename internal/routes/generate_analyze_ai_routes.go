package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/cache"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/ratelimit"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupGenerateAnalyzeAIRoutes configures AI filtering-related routes
func SetupGenerateAnalyzeAIRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize dependencies for CurriculumUseCase
	cacheService := cache.NewCacheService(redis.GetClient(), logger)
	curriculumRepo := repositories.NewCurriculumRepository(db, logger)
	curriculumUseCase := usecases.NewCurriculumUseCase(curriculumRepo, cacheService, logger)

	generateAnalyzeAIUseCase, err := usecases.NewGenerateAnalyzeAIUseCase(curriculumUseCase)
	if err != nil {
		logger.Error("Failed to create Generate Analyze AI usecase", zap.Error(err))
		return
	}

	generateAnalyzeAIHandler := handlers.NewGenerateAnalyzeAIHandler(generateAnalyzeAIUseCase)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateAnalyze := router.Group("/api/v1/generate-analyze-ai")
	generateAnalyze.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	generateAnalyze.Use(ratelimit.RateLimiterMiddleware(aiRateLimiter))
	{
		generateAnalyze.POST("/:id", generateAnalyzeAIHandler.FilterContent)
	}
}
