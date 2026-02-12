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
func SetupGenerateAcademicAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config, authMiddleware gin.HandlerFunc, subscriptionUseCase usecases.SubscriptionUseCase) {
	generateAcademicAIUseCase, err := usecases.NewGenerateAcademicAIUseCase(cfg.OpenAI.APIKey)
	if err != nil {
		logger.Error("Failed to create Generate Academic AI usecase", zap.Error(err))
		return
	}

	generateAcademicAIHandler := handlers.NewGenerateAcademicAIHandler(generateAcademicAIUseCase, logger)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateAcademic := router.Group(
		"/api/v1/generate-academic-ai",
		authMiddleware,
		middleware.RequireSubscriptionPlan(subscriptionUseCase, redis.GetClient(), config.DefaultAIQuotaByPlan()),
		ratelimit.RateLimiterMiddleware(aiRateLimiter),
	)
	{
		generateAcademic.POST("", generateAcademicAIHandler.FilterContent)
	}
}
