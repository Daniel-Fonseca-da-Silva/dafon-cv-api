package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/ratelimit"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) error {
	// Initialize rate limiter with environment configuration
	rateLimiter := ratelimit.NewDefaultRateLimiter(redis.GetClient(), logger)

	// Apply rate limiting to all routes except health check
	router.Use(ratelimit.RateLimiterMiddleware(rateLimiter))

	// Health check handler
	healthHandler := handlers.NewHealthCheckHandler()

	// Health check endpoint (no rate limiting)
	router.GET("/health", healthHandler.HealthCheck)

	// Setup user routes
	SetupUserRoutes(router, db, logger, cfg)

	// Setup curriculum routes
	SetupCurriculumRoutes(router, db, logger, cfg)

	// Setup AI analysis routes
	SetupGenerateIntroAIRoutes(router, logger, cfg)
	// Setup generate courses AI routes
	SetupGenerateCoursesAIRoutes(router, logger, cfg)

	// Setup generate academic AI routes
	SetupGenerateAcademicAIRoutes(router, logger, cfg)

	// Setup generate task AI routes
	SetupGenerateTaskAIRoutes(router, logger, cfg)

	// Setup generate skill AI routes
	SetupGenerateSkillAIRoutes(router, logger, cfg)

	// Setup configuration routes
	SetupConfigurationRoutes(router, db, logger, cfg)

	// Setup authentication email routes
	SetupEmailRoutes(router, logger, cfg)

	// Setup generate analyze AI routes
	SetupGenerateAnalyzeAIRoutes(router, logger, cfg)

	// Setup generate translation AI routes
	SetupGenerateTranslationAIRoutes(router, logger, cfg)

	return nil
}
