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

// SetupCurriculumRoutes configures curriculum-related routes
func SetupCurriculumRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize curriculum dependencies
	curriculumRepo := repositories.NewCurriculumRepository(db)
	curriculumUseCase := usecases.NewCurriculumUseCase(curriculumRepo)

	// Initialize user dependencies for user verification
	userRepo := repositories.NewUserRepository(db)
	configurationRepo := repositories.NewConfigurationRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo)

	curriculumHandler := handlers.NewCurriculumHandler(curriculumUseCase, userUseCase)

	curriculums := router.Group("/api/v1/curriculums")
	curriculums.Use(middleware.StaticTokenMiddleware(cfg.App.StaticToken))

	{
		curriculums.POST("", curriculumHandler.CreateCurriculum)
		curriculums.GET("/:id", curriculumHandler.GetCurriculumByID)
	}
}
