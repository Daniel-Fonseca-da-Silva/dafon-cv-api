package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

// GenerateCoursesAIHandler handles HTTP requests for AI filtering operations
type GenerateCoursesAIHandler struct {
	generateCoursesAIUseCase usecases.GenerateCoursesAIUseCase
}

// NewGenerateCoursesAIHandler creates a new instance of GenerateCoursesAIHandler
func NewGenerateCoursesAIHandler(generateCoursesAIUseCase usecases.GenerateCoursesAIUseCase) *GenerateCoursesAIHandler {
	return &GenerateCoursesAIHandler{
		generateCoursesAIUseCase: generateCoursesAIUseCase,
	}
}

// FilterContent handles POST /generate-intro-ai request
func (h *GenerateCoursesAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateCoursesAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	response, err := h.generateCoursesAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
