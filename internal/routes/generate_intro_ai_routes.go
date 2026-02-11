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

// SetupGenerateIntroAIRoutes configures AI filtering-related routes
func SetupGenerateIntroAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config, authMiddleware gin.HandlerFunc, subscriptionUseCase usecases.SubscriptionUseCase) {
	generateIntroAIUseCase, err := usecases.NewGenerateIntroAIUseCase(cfg.OpenAI.APIKey)
	if err != nil {
		logger.Error("Failed to create Generate Intro AI usecase", zap.Error(err))
		return
	}

	generateIntroAIHandler := handlers.NewGenerateIntroAIHandler(generateIntroAIUseCase, logger)

	// Create stricter rate limiter for AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateIntros := router.Group(
		"/api/v1/generate-intro-ai",
		authMiddleware,
		middleware.RequireSubscriptionPlan(subscriptionUseCase, redis.GetClient(), config.DefaultAIQuotaByPlan()),
		ratelimit.RateLimiterMiddleware(aiRateLimiter),
	)

	{
		generateIntros.POST("", generateIntroAIHandler.FilterContent)
	}
}
