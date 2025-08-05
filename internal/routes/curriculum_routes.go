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
func SetupCurriculumRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig := config.NewJWTConfig(logger)

	curriculumRepo := repositories.NewCurriculumRepository(db)
	curriculumUseCase := usecases.NewCurriculumUseCase(curriculumRepo)
	curriculumHandler := handlers.NewCurriculumHandler(curriculumUseCase)

	curriculums := router.Group("/api/v1/curriculums")
	curriculums.Use(middleware.AuthMiddleware(jwtConfig))
	{
		curriculums.POST("", curriculumHandler.CreateCurriculum)
		curriculums.GET("/:id", curriculumHandler.GetCurriculumByID)
	}
}
