package main

import (
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/database"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.Connect(cfg, logger); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run database migrations
	if err := database.AutoMigrate(logger); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Set Gin mode
	gin.SetMode(cfg.Mode)

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, database.GetDB(), logger)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Server starting",
		zap.String("port", cfg.Port),
		zap.String("mode", cfg.Mode),
	)
	logger.Info("Health check available",
		zap.String("url", fmt.Sprintf("http://localhost%s/health", addr)),
	)

	if err := router.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
