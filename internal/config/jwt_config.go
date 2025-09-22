package config

import (
	"errors"
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
func NewJWTConfig(logger *zap.Logger) (*JWTConfig, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		logger.Error("JWT_SECRET_KEY environment variable is required but not set")
		return nil, errors.New("JWT_SECRET_KEY environment variable is required")
	}
	logger.Info("JWT_SECRET_KEY loaded from environment")

	durationStr := os.Getenv("JWT_DURATION")
	if durationStr == "" {
		logger.Error("JWT_DURATION environment variable is required but not set")
		return nil, errors.New("JWT_DURATION environment variable is required")
	}
	logger.Info("JWT_DURATION loaded from environment",
		zap.String("duration", durationStr))

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("Failed to parse JWT_DURATION",
			zap.String("invalid_duration", durationStr),
			zap.Error(err))
		return nil, errors.New("invalid JWT_DURATION format")
	}
	logger.Info("JWT duration parsed successfully",
		zap.Duration("duration", duration))

	logger.Info("JWT configuration loaded successfully",
		zap.Int("secret_key_length", len(secretKey)),
		zap.Duration("token_duration", duration))

	return &JWTConfig{
		SecretKey: secretKey,
		Duration:  duration,
	}, nil
}
