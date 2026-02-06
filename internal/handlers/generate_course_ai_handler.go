package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateCoursesAIHandler handles HTTP requests for AI filtering operations
type GenerateCoursesAIHandler struct {
	generateCoursesAIUseCase usecases.GenerateCoursesAIUseCase
	logger                  *zap.Logger
}

// NewGenerateCoursesAIHandler creates a new instance of GenerateCoursesAIHandler
func NewGenerateCoursesAIHandler(generateCoursesAIUseCase usecases.GenerateCoursesAIUseCase, logger *zap.Logger) *GenerateCoursesAIHandler {
	return &GenerateCoursesAIHandler{
		generateCoursesAIUseCase: generateCoursesAIUseCase,
		logger:                  logger,
	}
}

// FilterContent godoc
// @Summary      Generate courses content with AI
// @Description  Generates courses/training content from the given input
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateCoursesAIRequest   true  "Content to process"
// @Success      200   {object}  dto.GenerateCoursesAIResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-courses-ai [post]
// @Security     BearerAuth
func (h *GenerateCoursesAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateCoursesAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateCoursesAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		h.abortWithInternalServerError(c, "generate courses content", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateCoursesAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate courses AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
