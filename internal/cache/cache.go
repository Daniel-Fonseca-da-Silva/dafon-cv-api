package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// CacheService handles Redis caching operations
type CacheService struct {
	client *redis.Client
	logger *zap.Logger
}

// NewCacheService creates a new cache service instance
func NewCacheService(client *redis.Client, logger *zap.Logger) *CacheService {
	return &CacheService{
		client: client,
		logger: logger,
	}
}

// Set stores a value in cache with TTL
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		cs.logger.Error("Failed to marshal cache value", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	err = cs.client.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		cs.logger.Error("Failed to set cache value", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	cs.logger.Debug("Cache value set successfully", zap.String("key", key), zap.Duration("ttl", ttl))
	return nil
}

// Get retrieves a value from cache
func (cs *CacheService) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := cs.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			cs.logger.Debug("Cache miss", zap.String("key", key))
			return false, nil
		}
		cs.logger.Error("Failed to get cache value", zap.Error(err), zap.String("key", key))
		return false, fmt.Errorf("failed to get cache value: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		cs.logger.Error("Failed to unmarshal cache value", zap.Error(err), zap.String("key", key))
		return false, fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	cs.logger.Debug("Cache hit", zap.String("key", key))
	return true, nil
}

// Delete removes a value from cache
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	err := cs.client.Del(ctx, key).Err()
	if err != nil {
		cs.logger.Error("Failed to delete cache value", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to delete cache value: %w", err)
	}

	cs.logger.Debug("Cache value deleted successfully", zap.String("key", key))
	return nil
}

// DeletePattern removes all keys matching a pattern
func (cs *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := cs.client.Keys(ctx, pattern).Result()
	if err != nil {
		cs.logger.Error("Failed to get keys by pattern", zap.Error(err), zap.String("pattern", pattern))
		return fmt.Errorf("failed to get keys by pattern: %w", err)
	}

	if len(keys) == 0 {
		cs.logger.Debug("No keys found for pattern", zap.String("pattern", pattern))
		return nil
	}

	err = cs.client.Del(ctx, keys...).Err()
	if err != nil {
		cs.logger.Error("Failed to delete keys by pattern", zap.Error(err), zap.String("pattern", pattern))
		return fmt.Errorf("failed to delete keys by pattern: %w", err)
	}

	cs.logger.Debug("Cache keys deleted by pattern", zap.String("pattern", pattern), zap.Int("count", len(keys)))
	return nil
}

// GenerateUserCacheKey generates cache key for user
func GenerateUserCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

// GenerateCurriculumCacheKey generates cache key for curriculum
func GenerateCurriculumCacheKey(curriculumID string) string {
	return fmt.Sprintf("curriculum:%s", curriculumID)
}
