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
func SetupGenerateSkillAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config, authMiddleware gin.HandlerFunc, subscriptionUseCase usecases.SubscriptionUseCase) {
	generateSkillAIUseCase, err := usecases.NewGenerateSkillAIUseCase(cfg.OpenAI.APIKey)
	if err != nil {
		logger.Error("Failed to create Generate Skill AI usecase", zap.Error(err))
		return
	}

	generateSkillAIHandler := handlers.NewGenerateSkillAIHandler(generateSkillAIUseCase, logger)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateSkill := router.Group(
		"/api/v1/generate-skill-ai",
		authMiddleware,
		middleware.RequireSubscriptionPlan(subscriptionUseCase, redis.GetClient(), config.DefaultAIQuotaByPlan()),
		ratelimit.RateLimiterMiddleware(aiRateLimiter),
	)
	{
		generateSkill.POST("", generateSkillAIHandler.FilterContent)
	}
}
