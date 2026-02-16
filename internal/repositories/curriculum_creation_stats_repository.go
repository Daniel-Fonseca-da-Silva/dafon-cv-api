package repositories

import (
	"context"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CurriculumCreationStatsRepository defines the interface for curriculum creation stats.
type CurriculumCreationStatsRepository interface {
	IncrementCreationCount(ctx context.Context, userID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type curriculumCreationStatsRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewCurriculumCreationStatsRepository creates a new CurriculumCreationStatsRepository.
func NewCurriculumCreationStatsRepository(db *gorm.DB, logger *zap.Logger) CurriculumCreationStatsRepository {
	return &curriculumCreationStatsRepository{db: db, logger: logger}
}

// IncrementCreationCount upserts a row for the user: insert with TotalCreations=1 or increment.
func (r *curriculumCreationStatsRepository) IncrementCreationCount(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stats models.CurriculumCreationStats
		err := tx.Where("user_id = ?", userID).First(&stats).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			stats = models.CurriculumCreationStats{UserID: userID, TotalCreations: 1}
			return tx.Create(&stats).Error
		}
		return tx.Model(&stats).Update("total_creations", gorm.Expr("total_creations + ?", 1)).Error
	})
	if err != nil {
		r.logger.Error("Failed to increment curriculum creation count",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return fmt.Errorf("failed to increment curriculum creation count: %w", err)
	}
	return nil
}

// GetByUserID returns TotalCreations for the user, or 0 if no row exists.
func (r *curriculumCreationStatsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var stats models.CurriculumCreationStats
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		r.logger.Error("Failed to get curriculum creation stats by user ID",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return 0, fmt.Errorf("failed to get curriculum creation stats by user ID %s: %w", userID.String(), err)
	}
	return stats.TotalCreations, nil
}
