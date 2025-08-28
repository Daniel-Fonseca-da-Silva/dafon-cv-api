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

// SetupPDFModernRoutes configures PDF generation routes
func SetupPDFModernRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig := config.NewJWTConfig(logger)

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize repositories
	curriculumRepo := repositories.NewCurriculumRepository(db)

	// Initialize use cases
	curriculumUseCase := usecases.NewCurriculumUseCase(curriculumRepo)
	modernPDFUseCase := usecases.NewModernPDFUseCase(curriculumUseCase, logger)

	// Initialize handlers with worker pool configuration
	modernPDFHandler := handlers.NewModernPDFHandler(
		modernPDFUseCase,
		logger,
		cfg.WorkerPool.NumWorkers,
		cfg.WorkerPool.QueueSize,
	)

	// PDF routes group
	pdfGroup := router.Group("/api/v1/template_modern")
	pdfGroup.Use(middleware.AuthMiddleware(jwtConfig))
	{
		pdfGroup.POST("/:id", modernPDFHandler.CreateModernPDF)
		pdfGroup.GET("/status", modernPDFHandler.GetPoolStatus) // Endpoint for monitoring
	}
}
