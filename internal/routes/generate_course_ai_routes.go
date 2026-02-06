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

// SetupGenerateCoursesAIRoutes configures AI filtering-related routes
func SetupGenerateCoursesAIRoutes(router *gin.Engine, logger *zap.Logger, cfg *config.Config) {
	generateCoursesAIUseCase, err := usecases.NewGenerateCoursesAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Courses AI usecase", zap.Error(err))
		return
	}

	generateCoursesAIHandler := handlers.NewGenerateCoursesAIHandler(generateCoursesAIUseCase, logger)
	// Criar rate limiter mais estrito para AI routes
	aiRateLimiter := ratelimit.NewAIRateLimiter(redis.GetClient(), logger)

	generateCourses := router.Group(
		"/api/v1/generate-courses-ai",
		middleware.StaticTokenMiddleware(cfg.App.StaticToken),
		ratelimit.RateLimiterMiddleware(aiRateLimiter),
	)
	{
		generateCourses.POST("", generateCoursesAIHandler.FilterContent)
	}
}
