// @title           Dafon CV API
// @version         1.0
// @description     API for consultancy business and cv generation
// @termsOfService  http://swagger.io/terms/

// @contact.name   DafonCV Support
// @contact.url    https://www.dafoncv.online
// @contact.email  dafondeveloper@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey UserIDHeader
// @in header
// @name X-User-ID
package main

import (
	"fmt"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/database"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/routes"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/validators"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// Register global custom validators for Gin binding
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterCustomValidators(v)
	}

	// Disable automatic redirects
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Configure CORS
	allowedOrigins := []string{
		"http://localhost:3000",
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
