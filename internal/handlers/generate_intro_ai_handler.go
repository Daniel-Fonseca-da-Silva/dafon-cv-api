package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// GenerateIntroAIHandler handles HTTP requests for AI filtering operations
type GenerateIntroAIHandler struct {
	generateIntroAIUseCase usecases.GenerateIntroAIUseCase
}

// NewGenerateIntroAIHandler creates a new instance of GenerateIntroAIHandler
func NewGenerateIntroAIHandler(generateIntroAIUseCase usecases.GenerateIntroAIUseCase) *GenerateIntroAIHandler {
	return &GenerateIntroAIHandler{
		generateIntroAIUseCase: generateIntroAIUseCase,
	}
}

// FilterContent handles POST /generate-intro-ai request
func (h *GenerateIntroAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateIntroAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	response, err := h.generateIntroAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
