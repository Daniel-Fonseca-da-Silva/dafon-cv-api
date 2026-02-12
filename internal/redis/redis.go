package redis

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var client *redis.Client

// Connect establishes connection to Redis using centralized configuration
func Connect(cfg *config.Config, logger *zap.Logger) error {
	var redisHost, redisPort, redisUsername, redisPassword, redisDBStr string
	var redisDB int

	// Check if REDIS_PUBLIC_URL is provided (priority)
	if cfg.Redis.PublicURL != "" {
		logger.Info("Using REDIS_PUBLIC_URL for Redis connection")

		redisURL, err := url.Parse(cfg.Redis.PublicURL)
		if err != nil {
			logger.Error("Failed to parse REDIS_PUBLIC_URL", zap.Error(err))
			return errors.WrapError(err, "invalid REDIS_PUBLIC_URL")
		}

		// Extract host and port
		redisHost = redisURL.Hostname()
		redisPort = redisURL.Port()

		// Extract username and password from UserInfo
		if redisURL.User != nil {
			redisUsername = redisURL.User.Username()
			redisPassword, _ = redisURL.User.Password()
		}

		// Extract DB from path (e.g., redis://host:port/0)
		if redisURL.Path != "" && len(redisURL.Path) > 1 {
			redisDBStr = strings.TrimPrefix(redisURL.Path, "/")
		} else {
			redisDBStr = "0"
		}
	} else {
		// Fallback to individual environment variables
		logger.Info("Using individual Redis environment variables")
		redisHost = cfg.Redis.Host
		redisPort = cfg.Redis.Port
		redisUsername = cfg.Redis.Username
		redisPassword = cfg.Redis.Password
		redisDBStr = cfg.Redis.DB
	}

	// Convert DB string to int
	var err error
	redisDB, err = strconv.Atoi(redisDBStr)
	if err != nil {
		logger.Warn("Invalid Redis DB value, using default 0", zap.String("db", redisDBStr))
		redisDB = 0
	}

	// Create Redis client
	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	}

	// Set username if provided (required for Redis 6.0+ with ACL or cloud providers)
	if redisUsername != "" {
		redisOptions.Username = redisUsername
	}

	client = redis.NewClient(redisOptions)

	// Test Redis connection
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return errors.WrapError(err, "Redis connection failed")
	}

	logger.Info("Successfully connected to Redis",
		zap.String("host", redisHost),
		zap.String("port", redisPort),
		zap.Int("db", redisDB))

	return nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// HealthCheck verifies Redis connection health
func HealthCheck() error {
	if client == nil {
		return errors.ErrRedisClientNotInitialized
	}

	// Create context with timeout for health check
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test Redis connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return errors.WrapError(err, "Redis health check failed")
	}

	return nil
}

// GetRedisInfo returns Redis server information
func GetRedisInfo() (map[string]string, error) {
	if client == nil {
		return nil, errors.ErrRedisClientNotInitialized
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get Redis server info
	info, err := client.Info(ctx, "server", "memory", "clients").Result()
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get Redis info")
	}

	// Parse info into map
	infoMap := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				infoMap[parts[0]] = parts[1]
			}
		}
	}

	return infoMap, nil
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
