package repositories

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"gorm.io/gorm"
)

// CurriculumRepository defines the interface for curriculum data operations
type CurriculumRepository interface {
	Create(ctx context.Context, curriculum *models.Curriculums) error
}

// curriculumRepository implements CurriculumRepository interface
type curriculumRepository struct {
	db *gorm.DB
}

// NewCurriculumRepository creates a new instance of CurriculumRepository
func NewCurriculumRepository(db *gorm.DB) CurriculumRepository {
	return &curriculumRepository{db: db}
}

// Create creates a new curriculum in the database
func (cu *curriculumRepository) Create(ctx context.Context, curriculum *models.Curriculums) error {
	// Use transaction to ensure atomicity when creating curriculum with works
	return cu.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(curriculum).Error
	})
}
