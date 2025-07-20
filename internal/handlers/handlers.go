package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct{}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// HealthCheck returns the health status of the application
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Application is running",
		"service": "dafon-cv-api",
	})
}
