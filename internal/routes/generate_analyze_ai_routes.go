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

// SetupGenerateAnalyzeAIRoutes configures AI filtering-related routes
func SetupGenerateAnalyzeAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config, authMiddleware gin.HandlerFunc, subscriptionUseCase usecases.SubscriptionUseCase, curriculumUseCase usecases.CurriculumUseCase) {
	generateAnalyzeAIUseCase, err := usecases.NewGenerateAnalyzeAIUseCase(cfg.OpenAI.APIKey, curriculumUseCase)
	if err != nil {
		logger.Error("Failed to create Generate Analyze AI usecase", zap.Error(err))
		return
	}

	generateAnalyzeAIHandler := handlers.NewGenerateAnalyzeAIHandler(generateAnalyzeAIUseCase, logger)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateAnalyze := router.Group(
		"/api/v1/generate-analyze-ai",
		authMiddleware,
		middleware.RequireSubscriptionPlan(subscriptionUseCase, redis.GetClient(), config.DefaultAIQuotaByPlan()),
		ratelimit.RateLimiterMiddleware(aiRateLimiter),
	)
	{
		generateAnalyze.POST("/:id", generateAnalyzeAIHandler.FilterContent)
	}
}
