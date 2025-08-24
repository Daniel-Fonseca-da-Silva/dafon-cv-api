package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// Health check handler
	healthHandler := handlers.NewHealthCheckHandler()

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// Setup auth routes
	SetupAuthRoutes(router, db, logger)

	// Setup user routes
	SetupUserRoutes(router, db, logger)

	// Setup curriculum routes
	SetupCurriculumRoutes(router, db, logger)

	// Setup AI analysis routes
	SetupGenerateIntroAIRoutes(router, logger)
	// Setup generate courses AI routes
	SetupGenerateCoursesAIRoutes(router, logger)

	// Setup generate task AI routes
	SetupGenerateTaskAIRoutes(router, logger)

	// Setup PDF routes
	SetupPDFSimpleRoutes(router, db, logger)

	// Setup configuration routes
	SetupConfigurationRoutes(router, db, logger)

}
