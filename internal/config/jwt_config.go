package config

import (
	"os"
	"time"

	"go.uber.org/zap"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
	Duration  time.Duration
}

// NewJWTConfig creates a new JWT configuration
func NewJWTConfig(logger *zap.Logger) *JWTConfig {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
		logger.Warn("JWT_SECRET_KEY not found in environment, using default secret key")
	} else {
		logger.Info("JWT_SECRET_KEY loaded from environment")
	}

	durationStr := os.Getenv("JWT_DURATION")
	if durationStr == "" {
		durationStr = "3h"
		logger.Warn("JWT_DURATION not found in environment, using default duration",
			zap.String("default_duration", durationStr))
	} else {
		logger.Info("JWT_DURATION loaded from environment",
			zap.String("duration", durationStr))
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration = 3 * time.Hour
		logger.Error("Failed to parse JWT_DURATION, using fallback duration",
			zap.String("invalid_duration", durationStr),
			zap.String("fallback_duration", "3h"),
			zap.Error(err))
	} else {
		logger.Info("JWT duration parsed successfully",
			zap.Duration("duration", duration))
	}

	logger.Info("JWT configuration loaded successfully",
		zap.Int("secret_key_length", len(secretKey)),
		zap.Duration("token_duration", duration))

	return &JWTConfig{
		SecretKey: secretKey,
		Duration:  duration,
	}
}
