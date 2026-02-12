package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GenerateAnalyzeAIHandler handles HTTP requests for AI filtering operations
type GenerateAnalyzeAIHandler struct {
	generateAnalyzeAIUseCase usecases.GenerateAnalyzeAIUseCase
	logger                   *zap.Logger
}

// NewGenerateAnalyzeAIHandler creates a new instance of GenerateAnalyzeAIHandler
func NewGenerateAnalyzeAIHandler(generateAnalyzeAIUseCase usecases.GenerateAnalyzeAIUseCase, logger *zap.Logger) *GenerateAnalyzeAIHandler {
	return &GenerateAnalyzeAIHandler{
		generateAnalyzeAIUseCase: generateAnalyzeAIUseCase,
		logger:                   logger,
	}
}

// FilterContent godoc
// @Summary      Analyze curriculum with AI
// @Description  Analyzes curriculum content and returns score, improvement points, ATS compatibility, and recommendations
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateAnalyzeAIRequest   true  "Curriculum content to analyze"
// @Success      200   {object}  dto.GenerateAnalyzeAIResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-analyze-ai [post]
// @Security     BearerAuth
func (h *GenerateAnalyzeAIHandler) FilterContent(c *gin.Context) {
	idParam := c.Param("id")
	curriculumID, err := uuid.Parse(idParam)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.NewAppError("invalid curriculum id"))
		return
	}

	// Get language parameter from query string (default: "pt")
	language := c.DefaultQuery("lang", "pt")

	// Validate language code
	validLanguages := map[string]bool{
		"pt": true,
		"en": true,
		"es": true,
	}
	if !validLanguages[language] {
		transporthttp.HandleValidationError(c, errors.NewAppError("invalid language code. Supported: pt, en, es"))
		return
	}

	aiResponse, err := h.generateAnalyzeAIUseCase.FilterContent(c.Request.Context(), curriculumID, language)
	if err != nil {
		h.abortWithInternalServerError(c, "analyze curriculum", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateAnalyzeAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate analyze AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
