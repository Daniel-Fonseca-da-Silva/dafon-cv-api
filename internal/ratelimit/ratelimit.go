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

// RateLimiter lida com o rate limiting usando Redis
type RateLimiter struct {
	client  *redis.Client
	limit   int
	windows time.Duration
	context context.Context
	logger  *zap.Logger
}

// getEnvAsInt obtém uma variável de ambiente como um inteiro com um valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// NewRateLimiter cria uma nova instância de rate limiter
func NewRateLimiter(client *redis.Client, limit int, windows time.Duration, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		client:  client,
		limit:   limit,
		windows: windows,
		context: context.Background(),
		logger:  logger,
	}
}

// NewDefaultRateLimiter cria um rate limiter com a configuração padrão do ambiente
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

// NewAIRateLimiter cria um rate limiter com limites mais rigorosos para endpoints AI
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

// Allow verifica se a solicitação é permitida com base na chave
func (rl *RateLimiter) Allow(key string) bool {
	// Incrementa o contador do Redis
	pipe := rl.client.TxPipeline()

	incr := pipe.Incr(rl.context, key)

	_, err := pipe.Exec(rl.context)
	if err != nil {
		rl.logger.Error("Failed to execute Redis pipeline", zap.Error(err))
		return false
	}

	// Se o contador for 1, define o tempo de expiração para a chave
	if incr.Val() == 1 {
		rl.client.Expire(rl.context, key, rl.windows)
	}

	allowed := incr.Val() <= int64(rl.limit)

	if !allowed {
		rl.logger.Warn("Rate limit exceeded",
			zap.Int64("current_count", incr.Val()),
			zap.Int("limit", rl.limit))
	}

	return allowed
}

// GetClientIP extrai o IP do cliente da solicitação
func GetClientIP(r *http.Request) string {
	// Verifica o cabeçalho X-Forwarded-For (para balanceadores de carga/proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
	}

	// Verifica o cabeçalho X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip, _, err := net.SplitHostPort(xri); err == nil {
			return ip
		}
	}

	// Caso não encontre, usa o RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// RateLimiterMiddleware cria um middleware para rate limiting
func RateLimiterMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipRateLimit(c.Request.URL.Path) {
			c.Next()
			return
		}

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

// UserRateLimiterMiddleware cria um middleware para rate limiting específico para usuários
func UserRateLimiterMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipRateLimit(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Tenta obter o ID do usuário do contexto (definido pelo middleware de autenticação)
		userID, exists := c.Get("user_id")
		if !exists {
			// Caso não encontre, usa o rate limiting baseado em IP
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
			// Usa o ID do usuário para rate limiting
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

func shouldSkipRateLimit(path string) bool {
	switch path {
	case "/health":
		return true
	case "/api/v1/subscriptions/webhook":
		return true
	default:
		return false
	}
}
