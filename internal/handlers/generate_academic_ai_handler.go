package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

// GenerateAcademicAIHandler handles HTTP requests for AI filtering operations
type GenerateAcademicAIHandler struct {
	generateAcademicAIUseCase usecases.GenerateAcademicAIUseCase
}

// NewGenerateAcademicAIHandler creates a new instance of GenerateAcademicAIHandler
func NewGenerateAcademicAIHandler(generateAcademicAIUseCase usecases.GenerateAcademicAIUseCase) *GenerateAcademicAIHandler {
	return &GenerateAcademicAIHandler{
		generateAcademicAIUseCase: generateAcademicAIUseCase,
	}
}

// FilterContent handles POST /generate-academic-ai request
func (h *GenerateAcademicAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateAcademicAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateAcademicAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
