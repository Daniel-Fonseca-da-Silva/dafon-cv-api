package repositories

import (
	"context"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CurriculumRepository defines the interface for curriculum data operations.
type CurriculumRepository interface {
	Create(ctx context.Context, curriculum *models.Curriculums) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Curriculums, error)
	GetPageAfterID(ctx context.Context, afterID *uuid.UUID, limit int) ([]models.Curriculums, bool, error)
	GetPageAfterIDByUserID(ctx context.Context, userID uuid.UUID, afterID *uuid.UUID, limit int) ([]models.Curriculums, bool, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Curriculums, error)
	Count(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteCurriculum(ctx context.Context, id uuid.UUID) error
}

// curriculumRepository implements CurriculumRepository.
type curriculumRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewCurriculumRepository creates a new CurriculumRepository.
func NewCurriculumRepository(db *gorm.DB, logger *zap.Logger) CurriculumRepository {
	return &curriculumRepository{db: db, logger: logger}
}

// Create persists a new curriculum and its relations atomically.
func (cu *curriculumRepository) Create(ctx context.Context, curriculum *models.Curriculums) error {
	err := cu.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(curriculum).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		cu.logger.Error("Failed to create curriculum",
			zap.Error(err),
			zap.String("curriculum_id", curriculum.ID.String()),
		)
		return fmt.Errorf("failed to create curriculum: %w", err)
	}
	return nil
}

// GetByID retrieves a curriculum by ID.
func (cu *curriculumRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Curriculums, error) {
	var curriculum models.Curriculums
	err := cu.db.WithContext(ctx).Preload("Works").Preload("Educations").Where("id = ?", id).First(&curriculum).Error
	if err != nil {
		cu.logger.Error("Failed to get curriculum by ID",
			zap.Error(err),
			zap.String("curriculum_id", id.String()),
		)
		return nil, fmt.Errorf("failed to get curriculum by ID %s: %w", id.String(), err)
	}
	return &curriculum, nil
}

// GetPageAfterID retrieves curriculums using cursor-based pagination.
// It orders by ID ascending and returns at most limit curriculums.
// The returned boolean indicates whether there is a next page.
func (cu *curriculumRepository) GetPageAfterID(ctx context.Context, afterID *uuid.UUID, limit int) ([]models.Curriculums, bool, error) {
	var curriculums []models.Curriculums

	if limit < 1 {
		return []models.Curriculums{}, false, nil
	}

	query := cu.db.WithContext(ctx).
		Preload("Works").
		Preload("Educations").
		Model(&models.Curriculums{})

	if afterID != nil && *afterID != uuid.Nil {
		query = query.Where("id > ?", afterID.String())
	}

	err := query.Order("id ASC").Limit(limit + 1).Find(&curriculums).Error
	if err != nil {
		cu.logger.Error("Failed to get curriculums with cursor pagination", zap.Error(err))
		return nil, false, fmt.Errorf("failed to get curriculums: %w", err)
	}

	hasNextPage := len(curriculums) > limit
	if hasNextPage {
		curriculums = curriculums[:limit]
	}

	return curriculums, hasNextPage, nil
}

// GetPageAfterIDByUserID retrieves curriculums for a given user using cursor-based pagination.
// It orders by ID ascending and returns at most limit curriculums.
// The returned boolean indicates whether there is a next page.
func (cu *curriculumRepository) GetPageAfterIDByUserID(ctx context.Context, userID uuid.UUID, afterID *uuid.UUID, limit int) ([]models.Curriculums, bool, error) {
	var curriculums []models.Curriculums

	if limit < 1 {
		return []models.Curriculums{}, false, nil
	}

	query := cu.db.WithContext(ctx).
		Preload("Works").
		Preload("Educations").
		Model(&models.Curriculums{}).
		Where("user_id = ?", userID)

	if afterID != nil && *afterID != uuid.Nil {
		query = query.Where("id > ?", afterID.String())
	}

	err := query.Order("id ASC").Limit(limit + 1).Find(&curriculums).Error
	if err != nil {
		cu.logger.Error("Failed to get curriculums by user with cursor pagination",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, false, fmt.Errorf("failed to get curriculums by user %s: %w", userID.String(), err)
	}

	hasNextPage := len(curriculums) > limit
	if hasNextPage {
		curriculums = curriculums[:limit]
	}

	return curriculums, hasNextPage, nil
}

// GetByUserID retrieves a curriculum by user ID.
func (cu *curriculumRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Curriculums, error) {
	var curriculum models.Curriculums
	err := cu.db.WithContext(ctx).Preload("Works").Preload("Educations").Where("user_id = ?", userID).First(&curriculum).Error
	if err != nil {
		cu.logger.Error("Failed to get curriculum by user ID",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to get curriculum by user ID %s: %w", userID.String(), err)
	}
	return &curriculum, nil
}

// Count returns the total number of curriculums.
func (cu *curriculumRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := cu.db.WithContext(ctx).Model(&models.Curriculums{}).Count(&count).Error
	if err != nil {
		cu.logger.Error("Failed to count curriculums", zap.Error(err))
		return 0, fmt.Errorf("failed to count curriculums: %w", err)
	}
	return count, nil
}

// CountByUserID returns the number of curriculums for the given user.
func (cu *curriculumRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := cu.db.WithContext(ctx).Model(&models.Curriculums{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		cu.logger.Error("Failed to count curriculums by user ID",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return 0, fmt.Errorf("failed to count curriculums by user ID %s: %w", userID.String(), err)
	}
	return count, nil
}

// DeleteCurriculum deletes a curriculum by ID.
func (cu *curriculumRepository) DeleteCurriculum(ctx context.Context, id uuid.UUID) error {
	if err := cu.db.WithContext(ctx).Delete(&models.Curriculums{}, id).Error; err != nil {
		cu.logger.Error("Failed to delete curriculum",
			zap.Error(err),
			zap.String("curriculum_id", id.String()),
		)
		return fmt.Errorf("failed to delete curriculum %s: %w", id.String(), err)
	}
	return nil
}
