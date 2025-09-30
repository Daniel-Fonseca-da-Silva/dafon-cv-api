package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthEmailHandler handles HTTP requests for authentication email operations
type AuthEmailHandler struct {
	emailUseCase usecases.EmailUseCase
	logger       *zap.Logger
}

// NewAuthEmailHandler creates a new instance of AuthEmailHandler
func NewAuthEmailHandler(emailUseCase usecases.EmailUseCase, logger *zap.Logger) *AuthEmailHandler {
	return &AuthEmailHandler{
		emailUseCase: emailUseCase,
		logger:       logger,
	}
}

// SendAuthEmail handles POST /auth/send-email request
func (h *AuthEmailHandler) SendAuthEmail(c *gin.Context) {
	var req dto.SendAuthEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	h.logger.Info("Processing authentication email request",
		zap.String("email", req.Email),
		zap.String("name", req.Name),
	)

	// Get base URL from request or use default
	baseURL := c.GetHeader("X-Forwarded-Proto") + "://" + c.GetHeader("X-Forwarded-Host")
	if baseURL == "://" {
		baseURL = "http://localhost:8080" // Default fallback
	}

	// Send the authentication email
	err := h.emailUseCase.SendSessionTokenEmail(req.Email, req.Name, req.URLToken, baseURL)
	if err != nil {
		h.logger.Error("Failed to send authentication email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		utils.HandleValidationError(c, err)
		return
	}

	h.logger.Info("Authentication email sent successfully",
		zap.String("email", req.Email),
	)

	response := dto.SendAuthEmailResponse{
		Message: "Authentication email sent successfully",
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
