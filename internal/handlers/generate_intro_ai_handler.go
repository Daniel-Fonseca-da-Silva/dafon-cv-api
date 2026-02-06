package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateIntroAIHandler handles HTTP requests for AI filtering operations
type GenerateIntroAIHandler struct {
	generateIntroAIUseCase usecases.GenerateIntroAIUseCase
	logger               *zap.Logger
}

// NewGenerateIntroAIHandler creates a new instance of GenerateIntroAIHandler
func NewGenerateIntroAIHandler(generateIntroAIUseCase usecases.GenerateIntroAIUseCase, logger *zap.Logger) *GenerateIntroAIHandler {
	return &GenerateIntroAIHandler{
		generateIntroAIUseCase: generateIntroAIUseCase,
		logger:               logger,
	}
}

// FilterContent godoc
// @Summary      Generate intro content with AI
// @Description  Generates professional introduction content from the given input
// @Tags         Generate AI
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateIntroAIRequest   true  "Content to process"
// @Success      200   {object}  dto.GenerateIntroAIResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/generate-intro-ai [post]
// @Security     BearerAuth
func (h *GenerateIntroAIHandler) FilterContent(c *gin.Context) {
	var req dto.GenerateIntroAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	aiResponse, err := h.generateIntroAIUseCase.FilterContent(c.Request.Context(), &req)
	if err != nil {
		h.abortWithInternalServerError(c, "generate intro content", err)
		return
	}

	c.JSON(http.StatusOK, aiResponse)
}

func (h *GenerateIntroAIHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Generate intro AI handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
