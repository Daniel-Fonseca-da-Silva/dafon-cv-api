package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var client *redis.Client

// Connect establishes connection to Redis using centralized configuration
func Connect(cfg *config.Config, logger *zap.Logger) error {
	// Get configuration from centralized config
	redisHost := cfg.Redis.Host
	redisPort := cfg.Redis.Port
	redisPassword := cfg.Redis.Password
	redisDBStr := cfg.Redis.DB

	// Set defaults if not configured
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	if redisDBStr == "" {
		redisDBStr = "0"
	}

	// Convert DB string to int
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		logger.Warn("Invalid Redis DB value, using default 0", zap.String("db", redisDBStr))
		redisDB = 0
	}

	// Create Redis client
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return err
	}

	logger.Info("Successfully connected to Redis")

	return nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// HealthCheck verifies Redis connection health
func HealthCheck() error {
	if client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	// Create context with timeout for health check
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test Redis connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

// GetRedisInfo returns Redis server information
func GetRedisInfo() (map[string]string, error) {
	if client == nil {
		return nil, fmt.Errorf("Redis client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get Redis server info
	info, err := client.Info(ctx, "server", "memory", "clients").Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to get Redis info: %w", err)
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
