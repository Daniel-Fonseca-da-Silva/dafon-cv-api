package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateTranslationAIHandler handles HTTP requests for AI filtering operations
type GenerateTranslationAIHandler struct {
	generateTranslationAIUseCase usecases.GenerateTranslationAIUseCase
	logger                     *zap.Logger
}

// NewGenerateTranslationAIHandler creates a new instance of GenerateTranslationAIHandler
func NewGenerateTranslationAIHandler(generateTranslationAIUseCase usecases.GenerateTranslationAIUseCase, logger *zap.Logger) *GenerateTranslationAIHandler {
	return &GenerateTranslationAIHandler{
		generateTranslationAIUseCase: generateTranslationAIUseCase,
		logger:                     logger,
	}
}

// FilterContent godoc
// @Summary      Translate curriculum with AI
// @Description  Translates curriculum content to the target language (pt, en, es)
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateTranslationAIRequestDoc  true  "Curriculum to translate + target_language (pt, en, es)"
// @Success      200   {object}  dto.CurriculumResponse  "Translated curriculum (same structure as create curriculum)"
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-translation-ai [post]
// @Security     BearerAuth
func (h *GenerateTranslationAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateTranslationAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateTranslationAIUseCase.TranslateCurriculum(c.Request.Context(), &req)
	if err != nil {
		h.abortWithInternalServerError(c, "translate curriculum", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateTranslationAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate translation AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
