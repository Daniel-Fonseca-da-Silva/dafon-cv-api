package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
	"github.com/gin-gonic/gin"
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

// FilterContent handles POST /generate-analyze-ai request for curriculum analysis
func (h *GenerateAnalyzeAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateAnalyzeAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	response, err := h.generateAnalyzeAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
