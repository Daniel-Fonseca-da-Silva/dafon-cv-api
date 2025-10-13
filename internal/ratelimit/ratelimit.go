package ratelimit

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimiter handles rate limiting using Redis
type RateLimiter struct {
	client  *redis.Client
	limit   int
	windows time.Duration
	context context.Context
	logger  *zap.Logger
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(client *redis.Client, limit int, windows time.Duration, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		client:  client,
		limit:   limit,
		windows: windows,
		context: context.Background(),
		logger:  logger,
	}
}

// NewDefaultRateLimiter creates a rate limiter with default configuration from environment
func NewDefaultRateLimiter(client *redis.Client, logger *zap.Logger) *RateLimiter {
	limit := getEnvAsInt("RATE_LIMIT", 100)
	windowMinutes := getEnvAsInt("RATE_WINDOW_MINUTES", 1)

	return &RateLimiter{
		client:  client,
		limit:   limit,
		windows: time.Duration(windowMinutes) * time.Minute,
		context: context.Background(),
		logger:  logger,
	}
}

// NewAIRateLimiter creates a rate limiter with stricter limits for AI endpoints
func NewAIRateLimiter(client *redis.Client, logger *zap.Logger) *RateLimiter {
	limit := getEnvAsInt("AI_RATE_LIMIT", 10)
	windowMinutes := getEnvAsInt("AI_RATE_WINDOW_MINUTES", 1)

	return &RateLimiter{
		client:  client,
		limit:   limit,
		windows: time.Duration(windowMinutes) * time.Minute,
		context: context.Background(),
		logger:  logger,
	}
}

// Allow checks if the request is allowed based on the key
func (rl *RateLimiter) Allow(key string) bool {
	pipe := rl.client.TxPipeline()

	incr := pipe.Incr(rl.context, key)
	pipe.Expire(rl.context, key, rl.windows)

	_, err := pipe.Exec(rl.context)
	if err != nil {
		rl.logger.Error("Failed to execute Redis pipeline", zap.Error(err))
		return false
	}

	allowed := incr.Val() <= int64(rl.limit)

	if !allowed {
		rl.logger.Warn("Rate limit exceeded",
			zap.Int64("current_count", incr.Val()),
			zap.Int("limit", rl.limit))
	}

	return allowed
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (for load balancers/proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
	}

	// Check for X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip, _, err := net.SplitHostPort(xri); err == nil {
			return ip
		}
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// RateLimiterMiddleware creates a middleware for rate limiting
func RateLimiterMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := GetClientIP(c.Request)

		if !rl.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimiterMiddleware creates a middleware for user-specific rate limiting
func UserRateLimiterMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// Fallback to IP-based rate limiting
			clientIP := GetClientIP(c.Request)
			if !rl.Allow(clientIP) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "Too Many Requests",
					"message": "Rate limit exceeded. Please try again later.",
				})
				c.Abort()
				return
			}
		} else {
			// Use user ID for rate limiting
			userKey := "user:" + strconv.Itoa(userID.(int))
			if !rl.Allow(userKey) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "Too Many Requests",
					"message": "Rate limit exceeded. Please try again later.",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
