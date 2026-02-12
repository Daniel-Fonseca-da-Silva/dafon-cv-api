package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateAcademicAIHandler handles HTTP requests for AI filtering operations
type GenerateAcademicAIHandler struct {
	generateAcademicAIUseCase usecases.GenerateAcademicAIUseCase
	logger                   *zap.Logger
}

// NewGenerateAcademicAIHandler creates a new instance of GenerateAcademicAIHandler
func NewGenerateAcademicAIHandler(generateAcademicAIUseCase usecases.GenerateAcademicAIUseCase, logger *zap.Logger) *GenerateAcademicAIHandler {
	return &GenerateAcademicAIHandler{
		generateAcademicAIUseCase: generateAcademicAIUseCase,
		logger:                   logger,
	}
}

// FilterContent godoc
// @Summary      Generate academic content with AI
// @Description  Generates academic-style content from the given input
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateAcademicAIRequest   true  "Content to process"
// @Success      200   {object}  dto.GenerateAcademicAIResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-academic-ai [post]
// @Security     BearerAuth
func (h *GenerateAcademicAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateAcademicAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateAcademicAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		h.abortWithInternalServerError(c, "generate academic content", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateAcademicAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate academic AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
