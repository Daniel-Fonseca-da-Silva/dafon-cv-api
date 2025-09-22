package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGenerateAcademicAIRoutes configures AI filtering-related routes
func SetupGenerateAcademicAIRoutes(router *gin.Engine, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig, err := config.NewJWTConfig(logger)
	if err != nil {
		logger.Fatal("Failed to initialize JWT config", zap.Error(err))
	}

	generateAcademicAIUseCase, err := usecases.NewGenerateAcademicAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Academic AI usecase", zap.Error(err))
		return
	}

	generateAcademicAIHandler := handlers.NewGenerateAcademicAIHandler(generateAcademicAIUseCase)

	generateAcademic := router.Group("/api/v1/generate-academic-ai")
	generateAcademic.Use(middleware.AuthMiddleware(jwtConfig))
	{
		generateAcademic.POST("", generateAcademicAIHandler.FilterContent)
	}
}
