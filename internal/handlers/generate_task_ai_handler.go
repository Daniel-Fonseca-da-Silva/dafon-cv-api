package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

// GenerateTaskAIHandler handles HTTP requests for AI filtering operations
type GenerateTaskAIHandler struct {
	generateTaskAIUseCase usecases.GenerateTaskAIUseCase
}

// NewGenerateTaskAIHandler creates a new instance of GenerateTaskAIHandler
func NewGenerateTaskAIHandler(generateTaskAIUseCase usecases.GenerateTaskAIUseCase) *GenerateTaskAIHandler {
	return &GenerateTaskAIHandler{
		generateTaskAIUseCase: generateTaskAIUseCase,
	}
}

// FilterContent handles POST /generate-task-ai request
func (h *GenerateTaskAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateTaskAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateTaskAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
