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

// SetupUserRoutes configures user-related routes
func SetupUserRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize user dependencies
	userRepo := repositories.NewUserRepository(db)
	configurationRepo := repositories.NewConfigurationRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Protected user routes (require authentication)
	protectedUsers := router.Group("/api/v1/user")
	protectedUsers.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))

	{
		protectedUsers.POST("", userHandler.CreateUser)
		protectedUsers.GET("/all", userHandler.GetAllUsers)   // Get all users
		protectedUsers.GET("/:id", userHandler.GetUserByID)   // Get user by ID
		protectedUsers.PATCH("/:id", userHandler.UpdateUser)  // Update user
		protectedUsers.DELETE("/:id", userHandler.DeleteUser) // Delete user
	}
}
