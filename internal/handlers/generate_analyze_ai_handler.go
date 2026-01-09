package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GenerateAnalyzeAIHandler handles HTTP requests for AI filtering operations
type GenerateAnalyzeAIHandler struct {
	generateAnalyzeAIUseCase usecases.GenerateAnalyzeAIUseCase
}

// NewGenerateAnalyzeAIHandler creates a new instance of GenerateAnalyzeAIHandler
func NewGenerateAnalyzeAIHandler(generateAnalyzeAIUseCase usecases.GenerateAnalyzeAIUseCase) *GenerateAnalyzeAIHandler {
	return &GenerateAnalyzeAIHandler{
		generateAnalyzeAIUseCase: generateAnalyzeAIUseCase,
	}
}

// FilterContent handles POST /generate-analyze-ai/:id request for curriculum analysis
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
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
