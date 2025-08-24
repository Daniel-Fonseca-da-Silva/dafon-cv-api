package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
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

// GetConfigurationByUserID handles GET /configuration/user/:user_id request
func (h *ConfigurationHandler) GetConfigurationByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	configuration, err := h.configurationUseCase.GetConfigurationByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found for this user"})
		return
	}

	c.JSON(http.StatusOK, configuration)
}

// UpdateConfiguration handles PATCH /configuration/:id request
func (h *ConfigurationHandler) UpdateConfiguration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration ID format"})
		return
	}

	var req dto.UpdateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	configuration, err := h.configurationUseCase.UpdateConfiguration(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configuration)
}

// DeleteConfiguration handles DELETE /configuration/:id request
func (h *ConfigurationHandler) DeleteConfiguration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration ID format"})
		return
	}

	if err := h.configurationUseCase.DeleteConfiguration(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}
