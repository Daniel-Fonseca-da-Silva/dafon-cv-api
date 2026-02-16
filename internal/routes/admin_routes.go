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

// SetupAdminRoutes configures admin (back office) routes.
// Double protection: X-Static-Token (trusted client) then Authorization Bearer session token; user must be admin.
func SetupAdminRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config, sessionRepo repositories.SessionRepository) {
	userRepo := repositories.NewUserRepository(db, logger)
	curriculumRepo := repositories.NewCurriculumRepository(db, logger)
	adminUseCase := usecases.NewAdminUseCase(userRepo, curriculumRepo, logger)
	adminHandler := handlers.NewAdminHandler(adminUseCase, logger)

	admin := router.Group(
		"/api/v1/admin",
		middleware.StaticTokenHeaderMiddleware(cfg.App.StaticToken),
		middleware.SessionMiddleware(sessionRepo),
		middleware.AdminMiddleware(userRepo),
	)
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
