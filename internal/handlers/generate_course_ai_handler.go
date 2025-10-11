package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/response"
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
		response.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateCoursesAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		response.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
