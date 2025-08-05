package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	response, err := h.generateTaskAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
