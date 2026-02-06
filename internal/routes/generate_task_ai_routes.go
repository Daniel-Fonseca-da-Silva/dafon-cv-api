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

// SetupGenerateTaskAIRoutes configures AI filtering-related routes
func SetupGenerateTaskAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateTaskAIUseCase, err := usecases.NewGenerateTaskAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Task AI usecase", zap.Error(err))
		return
	}

	generateTaskAIHandler := handlers.NewGenerateTaskAIHandler(generateTaskAIUseCase, logger)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateTasks := router.Group(
		"/api/v1/generate-task-ai",
		middleware.StaticTokenMiddleware(cfg.App.StaticToken),
		ratelimit.RateLimiterMiddleware(aiRateLimiter),
	)
	{
		generateTasks.POST("", generateTaskAIHandler.FilterContent)
	}
}
