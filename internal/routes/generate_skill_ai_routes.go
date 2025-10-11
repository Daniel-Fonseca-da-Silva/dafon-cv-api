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

// SetupGenerateSkillAIRoutes configures AI skill generation-related routes
func SetupGenerateSkillAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateSkillAIUseCase, err := usecases.NewGenerateSkillAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Skill AI usecase", zap.Error(err))
		return
	}

	generateSkillAIHandler := handlers.NewGenerateSkillAIHandler(generateSkillAIUseCase)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateSkill := router.Group("/api/v1/generate-skill-ai")
	generateSkill.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	generateSkill.Use(ratelimit.RateLimiterMiddleware(aiRateLimiter))
	{
		generateSkill.POST("", generateSkillAIHandler.FilterContent)
	}
}
