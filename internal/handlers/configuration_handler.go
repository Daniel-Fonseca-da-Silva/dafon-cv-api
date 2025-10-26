package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConfigurationHandler struct {
	configurationUseCase usecases.ConfigurationUseCase
}

func NewConfigurationHandler(configurationUseCase usecases.ConfigurationUseCase) *ConfigurationHandler {
	return &ConfigurationHandler{
		configurationUseCase: configurationUseCase,
	}
}

// Retorna uma configuração por ID do usuário
func (h *ConfigurationHandler) GetConfigurationByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	configuration, err := h.configurationUseCase.GetConfigurationByUserID(c.Request.Context(), userID)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("configuration not found for this user"))
		return
	}

	c.JSON(http.StatusOK, configuration)
}

// Atualiza uma configuração por ID do usuário
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
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, configuration)
}

// Deleta uma configuração por ID do usuário
func (h *ConfigurationHandler) DeleteConfiguration(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	if err := h.configurationUseCase.DeleteConfiguration(c.Request.Context(), userID); err != nil {
		transporthttp.HandleValidationError(c, errors.New("failed to delete configuration"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}
