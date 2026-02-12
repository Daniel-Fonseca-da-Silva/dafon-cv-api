package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateSkillAIHandler handles HTTP requests for AI skill generation operations
type GenerateSkillAIHandler struct {
	generateSkillAIUseCase usecases.GenerateSkillAIUseCase
	logger               *zap.Logger
}

// NewGenerateSkillAIHandler creates a new instance of GenerateSkillAIHandler
func NewGenerateSkillAIHandler(generateSkillAIUseCase usecases.GenerateSkillAIUseCase, logger *zap.Logger) *GenerateSkillAIHandler {
	return &GenerateSkillAIHandler{
		generateSkillAIUseCase: generateSkillAIUseCase,
		logger:               logger,
	}
}

// FilterContent godoc
// @Summary      Generate skills with AI
// @Description  Generates related skills from the given input
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateSkillAIRequest   true  "Content to process"
// @Success      200   {object}  dto.GenerateSkillAIResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-skill-ai [post]
// @Security     BearerAuth
func (h *GenerateSkillAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateSkillAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateSkillAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		h.abortWithInternalServerError(c, "generate skills", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateSkillAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate skill AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
