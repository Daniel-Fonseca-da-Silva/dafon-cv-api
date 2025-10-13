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

// SendEmail handles POST /email/send-email request
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
		h.logger.Error("Failed to send email")
		transporthttp.HandleValidationError(c, err)
		return
	}

	h.logger.Info("Email sent successfully")

	response := dto.SendEmailResponse{
		Message: "Email sent successfully",
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
