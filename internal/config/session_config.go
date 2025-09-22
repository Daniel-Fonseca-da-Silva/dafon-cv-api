package config

import (
	"errors"
	"os"
	"time"

	"go.uber.org/zap"
)

// SessionConfig holds session configuration
type SessionConfig struct {
	Duration time.Duration
}

// NewSessionConfig creates a new session configuration
func NewSessionConfig(logger *zap.Logger) (*SessionConfig, error) {
	durationStr := os.Getenv("SESSION_TOKEN_DURATION")
	if durationStr == "" {
		logger.Error("SESSION_TOKEN_DURATION environment variable is required but not set")
		return nil, errors.New("SESSION_TOKEN_DURATION environment variable is required")
	}
	logger.Info("SESSION_TOKEN_DURATION loaded from environment",
		zap.String("duration", durationStr))

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("Failed to parse SESSION_TOKEN_DURATION",
			zap.String("invalid_duration", durationStr),
			zap.Error(err))
		return nil, errors.New("invalid SESSION_TOKEN_DURATION format")
	}
	logger.Info("Session duration parsed successfully",
		zap.Duration("duration", duration))

	logger.Info("Session configuration loaded successfully",
		zap.Duration("token_duration", duration))

	return &SessionConfig{
		Duration: duration,
	}, nil
}
