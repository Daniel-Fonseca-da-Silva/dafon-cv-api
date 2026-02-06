package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/cache"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupUserRoutes configures user-related routes
func SetupUserRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize cache service
	cacheService := cache.NewCacheService(redis.GetClient(), logger)

	// Initialize user dependencies
	userRepo := repositories.NewUserRepository(db, logger)
	configurationRepo := repositories.NewConfigurationRepository(db, logger)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo, cacheService, logger)
	userHandler := handlers.NewUserHandler(userUseCase, logger)

	// Public user routes (no authentication)
	publicUsers := router.Group("/api/v1/user")
	publicUsers.POST("", userHandler.CreateUser)

	// Protected user routes (require authentication)
	protectedUsers := router.Group("/api/v1/user", middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	{
		protectedUsers.GET("/all", userHandler.GetAllUsers)
		protectedUsers.GET("/:id", userHandler.GetUserByID)
		protectedUsers.PATCH("/:id", userHandler.UpdateUser)
		protectedUsers.DELETE("/:id", userHandler.DeleteUser)
	}
}
