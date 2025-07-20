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
func SetupUserRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig := config.NewJWTConfig(logger)

	// Initialize user dependencies
	userRepo := repositories.NewUserRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(userUseCase)

	// User routes group (protected with authentication)
	users := router.Group("/api/v1/users")
	users.Use(middleware.AuthMiddleware(jwtConfig))
	{
		// Removendo POST / - criação de usuários agora é feita via /auth/register
		users.GET("/", userHandler.GetAllUsers)      // Get all users
		users.GET("/:id", userHandler.GetUserByID)   // Get user by ID
		users.PATCH("/:id", userHandler.UpdateUser)  // Update user
		users.DELETE("/:id", userHandler.DeleteUser) // Delete user
	}
}
