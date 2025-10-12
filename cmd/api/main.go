package main

import (
	"fmt"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/database"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/routes"
	"github.com/gin-contrib/cors"
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

	// Connect to Redis
	if err := redis.Connect(cfg, logger); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redis.Close()

	// Set Gin mode
	gin.SetMode(cfg.Mode)

	// Create router
	router := gin.Default()

	// Disable automatic redirects
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Configure CORS
	allowedOrigins := []string{
		"http://localhost:5173", "http://localhost:3000",
		"http://127.0.0.1:5173", "http://127.0.0.1:3000",
		"http://localhost:8080", "http://127.0.0.1:8080",
	}

	// Add Railway domain if APP_URL is set
	if cfg.App.URL != "" {
		allowedOrigins = append(allowedOrigins, cfg.App.URL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	if err := routes.SetupRoutes(router, database.GetDB(), logger, cfg); err != nil {
		logger.Fatal("Failed to setup routes", zap.Error(err))
	}

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
