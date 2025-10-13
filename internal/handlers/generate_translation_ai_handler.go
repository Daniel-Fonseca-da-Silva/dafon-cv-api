package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

// GenerateTranslationAIHandler handles HTTP requests for AI filtering operations
type GenerateTranslationAIHandler struct {
	generateTranslationAIUseCase usecases.GenerateTranslationAIUseCase
}

// NewGenerateTranslationAIHandler creates a new instance of GenerateTranslationAIHandler
func NewGenerateTranslationAIHandler(generateTranslationAIUseCase usecases.GenerateTranslationAIUseCase) *GenerateTranslationAIHandler {
	return &GenerateTranslationAIHandler{
		generateTranslationAIUseCase: generateTranslationAIUseCase,
	}
}

// FilterContent handles POST /generate-translation-ai request for curriculum translation
func (h *GenerateTranslationAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateTranslationAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateTranslationAIUseCase.TranslateCurriculum(c.Request.Context(), &req)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
