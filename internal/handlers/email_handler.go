package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EmailHandler handles HTTP requests for authentication email operations
type EmailHandler struct {
	emailUseCase usecases.EmailUseCase
	logger       *zap.Logger
}

// NewEmailHandler creates a new instance of EmailHandler
func NewEmailHandler(emailUseCase usecases.EmailUseCase, logger *zap.Logger) *EmailHandler {
	return &EmailHandler{
		emailUseCase: emailUseCase,
		logger:       logger,
	}
}

// SendEmail godoc
// @Summary      Send authentication email
// @Description  Sends an authentication email with session token link
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        body  body      dto.SendEmailRequest   true  "Email payload"
// @Success      200   {object}  dto.SendEmailResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/send-email [post]
// @Security     BearerAuth
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req dto.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	h.logger.Info("Processing email request")

	// Send the authentication email
	err := h.emailUseCase.SendSessionTokenEmail(req.Email, req.Name, req.URLToken)
	if err != nil {
		h.abortWithInternalServerError(c, "send session token email", err)
		return
	}

	h.logger.Info("Email sent successfully")

	response := dto.SendEmailResponse{
		Message: "Email sent successfully",
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}

func (h *EmailHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Email handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
