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

func SetupConfigurationRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize cache service
	cacheService := cache.NewCacheService(redis.GetClient(), logger)

	// Initialize configuration dependencies
	configurationRepo := repositories.NewConfigurationRepository(db)
	configurationUseCase := usecases.NewConfigurationUseCase(configurationRepo, cacheService, logger)
	configurationHandler := handlers.NewConfigurationHandler(configurationUseCase)

	// Configuration routes group (protected with authentication)
	configuration := router.Group("/api/v1/configuration")
	configuration.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))

	{
		configuration.GET("/:user_id", configurationHandler.GetConfigurationByUserID)
		configuration.PATCH("/:user_id", configurationHandler.UpdateConfiguration)
		configuration.DELETE("/:user_id", configurationHandler.DeleteConfiguration)
	}
}
