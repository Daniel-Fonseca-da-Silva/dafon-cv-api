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
}
