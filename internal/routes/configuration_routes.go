package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupConfigurationRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize configuration dependencies
	configurationRepo := repositories.NewConfigurationRepository(db)
	configurationUseCase := usecases.NewConfigurationUseCase(configurationRepo)
	configurationHandler := handlers.NewConfigurationHandler(configurationUseCase)

	// Configuration routes group (protected with authentication)
	configuration := router.Group("/api/v1/configuration")
	configuration.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))

	{
		configuration.GET("/user/:user_id", configurationHandler.GetConfigurationByUserID)
		configuration.PATCH("/:id", configurationHandler.UpdateConfiguration)
		configuration.DELETE("/:id", configurationHandler.DeleteConfiguration)
	}
}
