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

// SetupPDFSimpleRoutes configures PDF generation routes
func SetupPDFSimpleRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// Initialize JWT configuration
	jwtConfig := config.NewJWTConfig(logger)

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize repositories
	curriculumRepo := repositories.NewCurriculumRepository(db)

	// Initialize use cases
	curriculumUseCase := usecases.NewCurriculumUseCase(curriculumRepo)
	simplePDFUseCase := usecases.NewSimplePDFUseCase(curriculumUseCase, logger)

	// Initialize handlers with worker pool configuration
	simplePDFHandler := handlers.NewSimplePDFHandler(
		simplePDFUseCase,
		logger,
		cfg.WorkerPool.NumWorkers,
		cfg.WorkerPool.QueueSize,
	)

	// PDF routes group
	pdfGroup := router.Group("/api/v1/template_simple")
	pdfGroup.Use(middleware.AuthMiddleware(jwtConfig))
	{
		pdfGroup.POST("/:id", simplePDFHandler.CreateSimplePDF)
		pdfGroup.GET("/status", simplePDFHandler.GetPoolStatus) // Endpoint for monitoring
	}
}
