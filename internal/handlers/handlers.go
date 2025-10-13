package handlers

import (
	"net/http"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct {
	logger *zap.Logger
}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler(logger *zap.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{
		logger: logger,
	}
}

// HealthCheck returns the health status of the application
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	startTime := time.Now()

	// Check Redis health
	redisHealth := h.checkRedisHealth()
	redisHealthy := redisHealth["healthy"].(bool)

	// Calculate response time
	responseTime := time.Since(startTime)

	// Determine overall health status
	overallStatus := "ok"
	httpStatus := http.StatusOK

	if !redisHealthy {
		overallStatus = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	// Prepare response
	response := gin.H{
		"status":        overallStatus,
		"message":       "Application health check",
		"service":       "dafon-cv-api",
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"response_time": responseTime.String(),
		"components": gin.H{
			"redis": redisHealth,
		},
	}

	// Add Redis info if healthy
	if redisHealthy {
		if redisInfo, err := redis.GetRedisInfo(); err == nil {
			response["redis_info"] = gin.H{
				"version":           redisInfo["redis_version"],
				"uptime":            redisInfo["uptime_in_seconds"],
				"memory_used":       redisInfo["used_memory_human"],
				"connected_clients": redisInfo["connected_clients"],
			}
		}
	}

	c.JSON(httpStatus, response)
}

// checkRedisHealth checks Redis connection health
func (h *HealthCheckHandler) checkRedisHealth() gin.H {
	startTime := time.Now()

	err := redis.HealthCheck()
	responseTime := time.Since(startTime)

	if err != nil {
		h.logger.Error("Redis health check failed",
			zap.Error(err),
			zap.Duration("response_time", responseTime))

		return gin.H{
			"healthy":       false,
			"error":         err.Error(),
			"response_time": responseTime.String(),
		}
	}

	h.logger.Debug("Redis health check successful",
		zap.Duration("response_time", responseTime))

	return gin.H{
		"healthy":       true,
		"response_time": responseTime.String(),
	}
}
