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

// Login handles POST /auth/login request with email
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

// LoginWithToken handles GET /auth/login-with-token request with token in query parameter
func (h *AuthHandler) LoginWithToken(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		utils.HandleValidationError(c, errors.New("token parameter is required"))
		return
	}

	response, err := h.authUseCase.LoginWithToken(c.Request.Context(), token)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Set JWT as HTTP-only cookie
	c.SetCookie(
		"auth_token",    // name
		*response.Token, // value
		3600,            // maxAge (1 hour)
		"/",             // path
		"",              // domain
		false,           // secure (set to true in production with HTTPS)
		true,            // httpOnly
	)

	// Return JSON response for API testing (Postman) and frontend
	c.JSON(http.StatusOK, gin.H{
		"message":    "Login successful",
		"user":       response.User,
		"token":      response.Token,
		"expires_at": response.ExpiresAt,
	})
}
