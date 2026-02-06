package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateTaskAIHandler handles HTTP requests for AI filtering operations
type GenerateTaskAIHandler struct {
	generateTaskAIUseCase usecases.GenerateTaskAIUseCase
	logger              *zap.Logger
}

// NewGenerateTaskAIHandler creates a new instance of GenerateTaskAIHandler
func NewGenerateTaskAIHandler(generateTaskAIUseCase usecases.GenerateTaskAIUseCase, logger *zap.Logger) *GenerateTaskAIHandler {
	return &GenerateTaskAIHandler{
		generateTaskAIUseCase: generateTaskAIUseCase,
		logger:              logger,
	}
}

// FilterContent godoc
// @Summary      Generate task content with AI
// @Description  Generates task/description content from the given input
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateTaskAIRequest   true  "Content to process"
// @Success      200   {object}  dto.GenerateTaskAIResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-task-ai [post]
// @Security     BearerAuth
func (h *GenerateTaskAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateTaskAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateTaskAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		h.abortWithInternalServerError(c, "generate tasks content", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateTaskAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate task AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
