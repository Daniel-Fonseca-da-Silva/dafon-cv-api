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

// SetupAdminRoutes configures admin (back office) routes. All routes require X-User-ID header and admin user.
func SetupAdminRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	userRepo := repositories.NewUserRepository(db, logger)
	curriculumRepo := repositories.NewCurriculumRepository(db, logger)
	adminUseCase := usecases.NewAdminUseCase(userRepo, curriculumRepo, logger)
	adminHandler := handlers.NewAdminHandler(adminUseCase)

	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))
	admin.Use(middleware.AdminMiddleware(userRepo))
	{
		admin.GET("/dashboard", adminHandler.GetDashboard)

		// Users: register specific paths before parametric :id
		admin.GET("/users/stats", adminHandler.GetUsersStats)
		admin.GET("/users", adminHandler.GetUsers)
		admin.GET("/users/:id/detail", adminHandler.GetUserDetail)
		admin.PATCH("/users/:id/toggle-admin", adminHandler.ToggleAdmin)

		// Curriculums: register stats before list
		admin.GET("/curriculums/stats", adminHandler.GetCurriculumsStats)
		admin.GET("/curriculums", adminHandler.GetCurriculums)
	}
}
