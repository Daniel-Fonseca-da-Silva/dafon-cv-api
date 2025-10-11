package repositories

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CurriculumRepository Define a interface para operações de dados de curriculum
type CurriculumRepository interface {
	Create(ctx context.Context, curriculum *models.Curriculums) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Curriculums, error)
	GetAll(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]models.Curriculums, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int, sortBy, sortOrder string) ([]models.Curriculums, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Curriculums, error)
	DeleteCurriculum(ctx context.Context, id uuid.UUID) error
}

// curriculumRepository Implementa a interface CurriculumRepository
type curriculumRepository struct {
	db *gorm.DB
}

// NewCurriculumRepository Cria uma nova instância de CurriculumRepository
func NewCurriculumRepository(db *gorm.DB) CurriculumRepository {
	return &curriculumRepository{db: db}
}

// Create Cria um novo curriculum no banco de dados
func (cu *curriculumRepository) Create(ctx context.Context, curriculum *models.Curriculums) error {
	// Usar transaction para garantir atomicidade quando criar curriculum com works
	return cu.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(curriculum).Error
	})
}

// GetByID Recupera um curriculum por ID
func (cu *curriculumRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Curriculums, error) {
	var curriculum models.Curriculums
	err := cu.db.WithContext(ctx).Preload("Works").Preload("Educations").Where("id = ?", id).First(&curriculum).Error
	return &curriculum, err
}

// GetAll Recupera todos os curriculums paginados com works e educations
func (cu *curriculumRepository) GetAll(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]models.Curriculums, error) {
	var curriculums []models.Curriculums

	// Calcular offset para paginação
	offset := (page - 1) * pageSize

	// Validar e definir campos de ordenação permitidos
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"full_name":  true,
		"email":      true,
	}

	// Definir campo de ordenação padrão se não especificado ou inválido
	if sortBy == "" || !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	// Validar direção de ordenação
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Construir string de ordenação
	orderClause := sortBy + " " + sortOrder

	// Buscar curriculums com paginação e preload das relações
	err := cu.db.WithContext(ctx).
		Preload("Works").
		Preload("Educations").
		Offset(offset).
		Limit(pageSize).
		Order(orderClause).
		Find(&curriculums).Error

	return curriculums, err
}

// GetAllByUserID Recupera todos os curriculums paginados de um usuário específico com works e educations
func (cu *curriculumRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int, sortBy, sortOrder string) ([]models.Curriculums, error) {
	var curriculums []models.Curriculums

	// Calcular offset para paginação
	offset := (page - 1) * pageSize

	// Validar e definir campos de ordenação permitidos
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"full_name":  true,
		"email":      true,
	}

	// Definir campo de ordenação padrão se não especificado ou inválido
	if sortBy == "" || !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	// Validar direção de ordenação
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Construir string de ordenação
	orderClause := sortBy + " " + sortOrder

	// Buscar curriculums do usuário específico com paginação e preload das relações
	err := cu.db.WithContext(ctx).
		Preload("Works").
		Preload("Educations").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order(orderClause).
		Find(&curriculums).Error

	return curriculums, err
}

// GetByUserID Recupera um curriculum por user_id
func (cu *curriculumRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Curriculums, error) {
	var curriculum models.Curriculums
	err := cu.db.WithContext(ctx).Preload("Works").Preload("Educations").Where("user_id = ?", userID).First(&curriculum).Error
	return &curriculum, err
}

// DeleteCurriculum Deleta um curriculum por ID
func (cu *curriculumRepository) DeleteCurriculum(ctx context.Context, id uuid.UUID) error {
	return cu.db.WithContext(ctx).Delete(&models.Curriculums{}, id).Error
}
