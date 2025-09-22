package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/middleware"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGenerateSkillAIRoutes configures AI skill generation-related routes
func SetupGenerateSkillAIRoutes(router *gin.Engine, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig, err := config.NewJWTConfig(logger)
	if err != nil {
		logger.Fatal("Failed to initialize JWT config", zap.Error(err))
	}

	generateSkillAIUseCase, err := usecases.NewGenerateSkillAIUseCase()
	if err != nil {
		logger.Error("Failed to create Generate Skill AI usecase", zap.Error(err))
		return
	}

	generateSkillAIHandler := handlers.NewGenerateSkillAIHandler(generateSkillAIUseCase)

	generateSkill := router.Group("/api/v1/generate-skill-ai")
	generateSkill.Use(middleware.AuthMiddleware(jwtConfig))
	{
		generateSkill.POST("", generateSkillAIHandler.FilterContent)
	}
}
