package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) error {
	// Health check handler
	healthHandler := handlers.NewHealthCheckHandler()

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// Setup auth routes
	if err := SetupAuthRoutes(router, db, logger, cfg); err != nil {
		return errors.WrapError(err, "failed to setup auth routes")
	}

	// Setup user routes
	SetupUserRoutes(router, db, logger)

	// Setup curriculum routes
	SetupCurriculumRoutes(router, db, logger)

	// Setup AI analysis routes
	SetupGenerateIntroAIRoutes(router, logger)
	// Setup generate courses AI routes
	SetupGenerateCoursesAIRoutes(router, logger)

	// Setup generate academic AI routes
	SetupGenerateAcademicAIRoutes(router, logger)

	// Setup generate task AI routes
	SetupGenerateTaskAIRoutes(router, logger)

	// Setup configuration routes
	SetupConfigurationRoutes(router, db, logger)

	return nil
}
