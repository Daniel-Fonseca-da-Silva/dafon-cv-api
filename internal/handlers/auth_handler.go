package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	authUseCase usecases.AuthUseCase
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authUseCase usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Login handles POST /auth/login request
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	response, err := h.authUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register handles POST /auth/register request
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Validate strong password
	if !utils.ValidatePasswordAndReturnError(c, req.Password) {
		return
	}

	response, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Logout handles POST /auth/logout request
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleValidationError(c, errors.New("user not authenticated"))
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		utils.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	err := h.authUseCase.Logout(c.Request.Context(), userIDStr)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	response := dto.LogoutResponse{
		Message: "Logout successful",
	}
	c.JSON(http.StatusOK, response)
}

// ForgotPassword handles POST /auth/forgot-password request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	response, err := h.authUseCase.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword handles POST /auth/reset-password?token=<token> request
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		utils.HandleValidationError(c, errors.New("token is required"))
		return
	}

	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Set token from query parameter
	req.Token = token

	// Validate strong password
	if !utils.ValidatePasswordAndReturnError(c, req.Password) {
		return
	}

	response, err := h.authUseCase.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
