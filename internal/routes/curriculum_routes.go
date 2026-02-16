package routes

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/cache"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/handlers"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupCurriculumRoutes configures curriculum-related routes
func SetupCurriculumRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, cfg *config.Config, authMiddleware gin.HandlerFunc, curriculumUseCase usecases.CurriculumUseCase) {
	// Initialize cache service (used by user use case)
	cacheService := cache.NewCacheService(redis.GetClient(), logger)

	// Initialize user dependencies for user verification
	userRepo := repositories.NewUserRepository(db, logger)
	configurationRepo := repositories.NewConfigurationRepository(db, logger)
	subscriptionRepo := repositories.NewSubscriptionRepository(db, logger)
	userUseCase := usecases.NewUserUseCase(userRepo, configurationRepo, subscriptionRepo, cacheService, logger)

	curriculumHandler := handlers.NewCurriculumHandler(curriculumUseCase, userUseCase, logger)

	curriculums := router.Group("/api/v1/curriculums", authMiddleware)

	{
		curriculums.POST("", curriculumHandler.CreateCurriculum)
		curriculums.GET("/get-all-by-user/:user_id", curriculumHandler.GetAllCurriculums)
		curriculums.GET("/count-by-user/:user_id", curriculumHandler.GetCurriculumCountByUserID)
		curriculums.GET("/creation-count-by-user/:user_id", curriculumHandler.GetCreationCountByUserID)
		curriculums.GET("/:curriculum_id", curriculumHandler.GetCurriculumByID)
		curriculums.GET("/get-body/:curriculum_id", curriculumHandler.GetCurriculumBody)
		curriculums.DELETE("/:curriculum_id", curriculumHandler.DeleteCurriculum)
	}
}
