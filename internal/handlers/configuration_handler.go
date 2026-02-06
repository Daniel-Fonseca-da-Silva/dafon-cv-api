package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ConfigurationHandler struct {
	configurationUseCase usecases.ConfigurationUseCase
	logger               *zap.Logger
}

func NewConfigurationHandler(configurationUseCase usecases.ConfigurationUseCase, logger *zap.Logger) *ConfigurationHandler {
	return &ConfigurationHandler{
		configurationUseCase: configurationUseCase,
		logger:               logger,
	}
}

// GetConfigurationByUserID godoc
// @Summary      Get configuration by user ID
// @Description  Returns the configuration for the given user ID
// @Tags         configuration
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {object}  dto.ConfigurationResponse
// @Failure      400      {object}  dto.ErrorResponseValidation  "Invalid user ID format"
// @Failure      404      {object}  dto.ErrorResponse  "Configuration not found"
// @Failure      500      {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/configuration/{user_id} [get]
// @Security     BearerAuth
func (h *ConfigurationHandler) GetConfigurationByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	configuration, err := h.configurationUseCase.GetConfigurationByUserID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "configuration not found for this user")
			return
		}
		h.abortWithInternalServerError(c, "get configuration by user id", err)
		return
	}

	c.JSON(http.StatusOK, configuration)
}

// UpdateConfiguration godoc
// @Summary      Update configuration by user ID
// @Description  Updates the configuration for the given user ID
// @Tags         configuration
// @Accept       json
// @Produce      json
// @Param        user_id  path      string                            true  "User ID"
// @Param        body     body      dto.UpdateConfigurationRequest   true  "Update payload"
// @Success      200      {object}  dto.ConfigurationResponse
// @Failure      400      {object}  dto.ErrorResponseValidation  "Invalid user ID format or validation error"
// @Failure      404      {object}  dto.ErrorResponse  "Configuration not found"
// @Failure      500      {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/configuration/{user_id} [patch]
// @Security     BearerAuth
func (h *ConfigurationHandler) UpdateConfiguration(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	var req dto.UpdateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	configuration, err := h.configurationUseCase.UpdateConfiguration(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "configuration not found for this user")
			return
		}
		h.abortWithInternalServerError(c, "update configuration", err)
		return
	}

	c.JSON(http.StatusOK, configuration)
}

// DeleteConfiguration godoc
// @Summary      Delete configuration by user ID
// @Description  Deletes the configuration for the given user ID
// @Tags         configuration
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  dto.ErrorResponseValidation  "Invalid user ID format"
// @Failure      404      {object}  dto.ErrorResponse  "Configuration not found"
// @Failure      500      {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/configuration/{user_id} [delete]
// @Security     BearerAuth
func (h *ConfigurationHandler) DeleteConfiguration(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	if err := h.configurationUseCase.DeleteConfiguration(c.Request.Context(), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "configuration not found for this user")
			return
		}
		h.abortWithInternalServerError(c, "delete configuration", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}

func (h *ConfigurationHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Configuration handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
