package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

// GenerateSkillAIHandler handles HTTP requests for AI skill generation operations
type GenerateSkillAIHandler struct {
	generateSkillAIUseCase usecases.GenerateSkillAIUseCase
}

// NewGenerateSkillAIHandler creates a new instance of GenerateSkillAIHandler
func NewGenerateSkillAIHandler(generateSkillAIUseCase usecases.GenerateSkillAIUseCase) *GenerateSkillAIHandler {
	return &GenerateSkillAIHandler{
		generateSkillAIUseCase: generateSkillAIUseCase,
	}
}

// FilterContent handles POST /generate-skill-ai request to generate related skills
func (h *GenerateSkillAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateSkillAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateSkillAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}
