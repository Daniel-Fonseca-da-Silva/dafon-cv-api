package usecases

import (
	"context"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/cache"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CurriculumUsecase Define a interface para operações de dados de curriculum
type CurriculumUseCase interface {
	CreateCurriculum(ctx context.Context, userID uuid.UUID, req *dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error)
	GetCurriculumByID(ctx context.Context, id uuid.UUID) (*dto.CurriculumResponse, error)
	GetAllCurriculums(ctx context.Context, userID uuid.UUID, page, pageSize int, sortBy, sortOrder string) ([]dto.CurriculumResponse, error)
	GetCurriculumBody(ctx context.Context, curriculumID uuid.UUID) (*dto.CurriculumBodyResponse, error)
	DeleteCurriculum(ctx context.Context, id uuid.UUID) error
}

// curriculumUsecase Implementa a interface CurriculumUseCase
type curriculumUseCase struct {
	curriculumRepo repositories.CurriculumRepository
	cacheService   *cache.CacheService
	logger         *zap.Logger
}

// NewCurriculumUsecase Cria uma nova instância de CurriculumUsecase
func NewCurriculumUseCase(curriculumRepo repositories.CurriculumRepository, cacheService *cache.CacheService, logger *zap.Logger) CurriculumUseCase {
	return &curriculumUseCase{
		curriculumRepo: curriculumRepo,
		cacheService:   cacheService,
		logger:         logger,
	}
}

// CreateCurriculum Cria um novo curriculum no banco de dados
func (cu *curriculumUseCase) CreateCurriculum(ctx context.Context, userID uuid.UUID, req *dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error) {
	// Create curriculum model
	curriculum := &models.Curriculums{
		FullName:      req.FullName,
		Email:         req.Email,
		Phone:         req.Phone,
		DriverLicense: req.DriverLicense,
		Intro:         req.Intro,
		Skills:        req.Skills,
		Languages:     req.Languages,
		Courses:       req.Courses,
		SocialLinks:   req.SocialLinks,
		ImageURL:      req.ImageURL,
		UserID:        userID,
	}

	// Validar o modelo de curriculum
	validate := validator.New()
	// Register custom validators
	validators.RegisterCustomValidators(validate)
	if err := validate.Struct(curriculum); err != nil {
		return nil, err
	}

	// Criar works associados ao curriculum
	for _, workReq := range req.Works {
		work := models.Work{
			Position:    workReq.Position,
			Company:     workReq.Company,
			Description: workReq.Description,
			StartDate:   workReq.StartDate,
			EndDate:     workReq.EndDate,
		}
		curriculum.Works = append(curriculum.Works, work)
	}

	// Criar educations associados ao curriculum
	for _, educationReq := range req.Educations {
		education := models.Education{
			Institution: educationReq.Institution,
			Degree:      educationReq.Degree,
			StartDate:   educationReq.StartDate,
			EndDate:     educationReq.EndDate,
			Description: educationReq.Description,
		}
		curriculum.Educations = append(curriculum.Educations, education)
	}

	// Salvar no banco de dados (GORM irá lidar com a relação de chave estrangeira)
	if err := cu.curriculumRepo.Create(ctx, curriculum); err != nil {
		return nil, err
	}

	// Preparar response de works
	worksResponse := make([]dto.WorkResponse, 0, len(curriculum.Works))
	for _, work := range curriculum.Works {
		worksResponse = append(worksResponse, dto.WorkResponse{
			ID:          work.ID,
			Position:    work.Position,
			Company:     work.Company,
			Description: work.Description,
			StartDate:   work.StartDate,
			EndDate:     work.EndDate,
			CreatedAt:   work.CreatedAt,
			UpdatedAt:   work.UpdatedAt,
		})
	}

	// Preparar response de educations
	educationsResponse := make([]dto.EducationResponse, 0, len(curriculum.Educations))
	for _, education := range curriculum.Educations {
		educationsResponse = append(educationsResponse, dto.EducationResponse{
			ID:          education.ID,
			Institution: education.Institution,
			Degree:      education.Degree,
			StartDate:   education.StartDate,
			EndDate:     education.EndDate,
			Description: education.Description,
			CreatedAt:   education.CreatedAt,
			UpdatedAt:   education.UpdatedAt,
		})
	}

	// Return response
	return &dto.CurriculumResponse{
		ID:            curriculum.ID,
		FullName:      curriculum.FullName,
		Email:         curriculum.Email,
		Phone:         curriculum.Phone,
		DriverLicense: curriculum.DriverLicense,
		Intro:         curriculum.Intro,
		Skills:        curriculum.Skills,
		Languages:     curriculum.Languages,
		Courses:       curriculum.Courses,
		SocialLinks:   curriculum.SocialLinks,
		ImageURL:      curriculum.ImageURL,
		Works:         worksResponse,
		Educations:    educationsResponse,
		CreatedAt:     curriculum.CreatedAt,
		UpdatedAt:     curriculum.UpdatedAt,
	}, nil
}

// GetCurriculumByID retrieves a curriculum by ID
func (cu *curriculumUseCase) GetCurriculumByID(ctx context.Context, id uuid.UUID) (*dto.CurriculumResponse, error) {
	cacheKey := cache.GenerateCurriculumCacheKey(id.String())

	// Try to get from cache first
	var curriculumResponse dto.CurriculumResponse
	found, err := cu.cacheService.Get(ctx, cacheKey, &curriculumResponse)
	if err != nil {
		cu.logger.Warn("Failed to get curriculum from cache, falling back to database",
			zap.Error(err),
			zap.String("curriculum_id", id.String()))
	} else if found {
		cu.logger.Debug("Curriculum retrieved from cache", zap.String("curriculum_id", id.String()))
		return &curriculumResponse, nil
	}

	// Cache miss - get from database
	curriculum, err := cu.curriculumRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert works from models to DTOs
	worksResponse := make([]dto.WorkResponse, 0, len(curriculum.Works))
	for _, work := range curriculum.Works {
		worksResponse = append(worksResponse, dto.WorkResponse{
			ID:          work.ID,
			Position:    work.Position,
			Company:     work.Company,
			Description: work.Description,
			StartDate:   work.StartDate,
			EndDate:     work.EndDate,
			CreatedAt:   work.CreatedAt,
			UpdatedAt:   work.UpdatedAt,
		})
	}

	// Convert educations from models to DTOs
	educationsResponse := make([]dto.EducationResponse, 0, len(curriculum.Educations))
	for _, education := range curriculum.Educations {
		educationsResponse = append(educationsResponse, dto.EducationResponse{
			ID:          education.ID,
			Institution: education.Institution,
			Degree:      education.Degree,
			StartDate:   education.StartDate,
			EndDate:     education.EndDate,
			Description: education.Description,
			CreatedAt:   education.CreatedAt,
			UpdatedAt:   education.UpdatedAt,
		})
	}

	// Create response
	curriculumResponse = dto.CurriculumResponse{
		ID:            curriculum.ID,
		FullName:      curriculum.FullName,
		Email:         curriculum.Email,
		Phone:         curriculum.Phone,
		DriverLicense: curriculum.DriverLicense,
		Intro:         curriculum.Intro,
		Skills:        curriculum.Skills,
		Languages:     curriculum.Languages,
		Courses:       curriculum.Courses,
		SocialLinks:   curriculum.SocialLinks,
		ImageURL:      curriculum.ImageURL,
		Works:         worksResponse,
		Educations:    educationsResponse,
		CreatedAt:     curriculum.CreatedAt,
		UpdatedAt:     curriculum.UpdatedAt,
	}

	// Armazena os dados em cache por 30 minutos
	ttl := 30 * time.Minute
	if err := cu.cacheService.Set(ctx, cacheKey, curriculumResponse, ttl); err != nil {
		cu.logger.Warn("Failed to cache curriculum data",
			zap.Error(err),
			zap.String("curriculum_id", id.String()),
			zap.Duration("ttl", ttl))
	} else {
		cu.logger.Debug("Curriculum cached successfully",
			zap.String("curriculum_id", id.String()),
			zap.Duration("ttl", ttl))
	}

	return &curriculumResponse, nil
}

// GetAllCurriculums traz todos os curriculums paginados de um usuário específico
func (cu *curriculumUseCase) GetAllCurriculums(ctx context.Context, userID uuid.UUID, page, pageSize int, sortBy, sortOrder string) ([]dto.CurriculumResponse, error) {
	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Tamanho de página padrão
	}

	curriculums, err := cu.curriculumRepo.GetAllByUserID(ctx, userID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// Converter curriculums para DTOs
	curriculumsResponse := make([]dto.CurriculumResponse, 0, len(curriculums))
	for _, curriculum := range curriculums {
		// Converter works de models para DTOs
		worksResponse := make([]dto.WorkResponse, 0, len(curriculum.Works))
		for _, work := range curriculum.Works {
			worksResponse = append(worksResponse, dto.WorkResponse{
				ID:          work.ID,
				Position:    work.Position,
				Company:     work.Company,
				Description: work.Description,
				StartDate:   work.StartDate,
				EndDate:     work.EndDate,
				CreatedAt:   work.CreatedAt,
				UpdatedAt:   work.UpdatedAt,
			})
		}

		// Converter educations de models para DTOs
		educationsResponse := make([]dto.EducationResponse, 0, len(curriculum.Educations))
		for _, education := range curriculum.Educations {
			educationsResponse = append(educationsResponse, dto.EducationResponse{
				ID:          education.ID,
				Institution: education.Institution,
				Degree:      education.Degree,
				StartDate:   education.StartDate,
				EndDate:     education.EndDate,
				Description: education.Description,
				CreatedAt:   education.CreatedAt,
				UpdatedAt:   education.UpdatedAt,
			})
		}

		curriculumsResponse = append(curriculumsResponse, dto.CurriculumResponse{
			ID:            curriculum.ID,
			FullName:      curriculum.FullName,
			Email:         curriculum.Email,
			Phone:         curriculum.Phone,
			DriverLicense: curriculum.DriverLicense,
			Intro:         curriculum.Intro,
			Skills:        curriculum.Skills,
			Languages:     curriculum.Languages,
			Courses:       curriculum.Courses,
			SocialLinks:   curriculum.SocialLinks,
			ImageURL:      curriculum.ImageURL,
			Works:         worksResponse,
			Educations:    educationsResponse,
			CreatedAt:     curriculum.CreatedAt,
			UpdatedAt:     curriculum.UpdatedAt,
		})
	}

	return curriculumsResponse, nil
}

// GetCurriculumBody retrieves a curriculum body in text format by curriculum ID
func (cu *curriculumUseCase) GetCurriculumBody(ctx context.Context, curriculumID uuid.UUID) (*dto.CurriculumBodyResponse, error) {
	cacheKey := cache.GenerateCurriculumBodyCacheKey(curriculumID.String())

	// Try to get from cache first
	var curriculumBodyResponse dto.CurriculumBodyResponse
	found, err := cu.cacheService.Get(ctx, cacheKey, &curriculumBodyResponse)
	if err != nil {
		cu.logger.Warn("Failed to get curriculum body from cache, falling back to database",
			zap.Error(err),
			zap.String("curriculum_id", curriculumID.String()))
	} else if found {
		cu.logger.Debug("Curriculum body retrieved from cache", zap.String("curriculum_id", curriculumID.String()))
		return &curriculumBodyResponse, nil
	}

	// Cache miss - get from database
	curriculum, err := cu.curriculumRepo.GetByID(ctx, curriculumID)
	if err != nil {
		return nil, err
	}

	// Build curriculum body in text format
	body := buildCurriculumBodyText(curriculum)

	// Create response
	curriculumBodyResponse = dto.CurriculumBodyResponse{
		Body: body,
	}

	// Armazena os dados em cache por 30 minutos
	ttl := 30 * time.Minute
	if err := cu.cacheService.Set(ctx, cacheKey, curriculumBodyResponse, ttl); err != nil {
		cu.logger.Warn("Failed to cache curriculum body data",
			zap.Error(err),
			zap.String("curriculum_id", curriculumID.String()),
			zap.Duration("ttl", ttl))
	} else {
		cu.logger.Debug("Curriculum body cached successfully",
			zap.String("curriculum_id", curriculumID.String()),
			zap.Duration("ttl", ttl))
	}

	return &curriculumBodyResponse, nil
}

// buildCurriculumBodyText builds the curriculum body in plain text format
func buildCurriculumBodyText(curriculum *models.Curriculums) string {
	var body string

	// Personal Information
	body += "Personal Information "
	body += "Name: " + curriculum.FullName + " "
	body += "Email: " + curriculum.Email + " "
	body += "Phone: " + curriculum.Phone + " "
	body += "Driver License: " + curriculum.DriverLicense + " "

	// Introduction
	if curriculum.Intro != "" {
		body += "Presentation " + curriculum.Intro + " "
	}

	// Skills
	if curriculum.Skills != "" {
		body += "Skills " + curriculum.Skills + " "
	}

	// Languages
	if curriculum.Languages != "" {
		body += "Languages " + curriculum.Languages + " "
	}

	// Courses
	if curriculum.Courses != "" {
		body += "Courses " + curriculum.Courses + " "
	}

	// Social Links
	if curriculum.SocialLinks != "" {
		body += "Social Links " + curriculum.SocialLinks + " "
	}

	// Image URL
	if curriculum.ImageURL != nil && *curriculum.ImageURL != "" {
		body += "Image URL " + *curriculum.ImageURL + " "
	}

	// Work Experience
	if len(curriculum.Works) > 0 {
		body += "Work Experience "
		for _, work := range curriculum.Works {
			body += "Position: " + work.Position + " "
			body += "Company: " + work.Company + " "
			startDate := work.StartDate.Format("01/02/2006")
			endDate := "Current"
			if work.EndDate != nil {
				endDate = work.EndDate.Format("01/02/2006")
			}
			body += "Period: " + startDate + " - " + endDate + " "
			if work.Description != "" {
				body += "Description: " + work.Description + " "
			}
		}
	}

	// Education
	if len(curriculum.Educations) > 0 {
		body += "Academic Formation "
		for _, education := range curriculum.Educations {
			body += "Degree: " + education.Degree + " "
			body += "Institution: " + education.Institution + " "
			startDate := education.StartDate.Format("01/02/2006")
			endDate := "Current"
			if education.EndDate != nil {
				endDate = education.EndDate.Format("01/02/2006")
			}
			body += "Period: " + startDate + " - " + endDate + " "
			if education.Description != "" {
				body += "Description: " + education.Description + " "
			}
		}
	}

	return body
}

// DeleteCurriculum Deleta um curriculum por ID
func (cu *curriculumUseCase) DeleteCurriculum(ctx context.Context, id uuid.UUID) error {
	// Delete curriculum from database
	if err := cu.curriculumRepo.DeleteCurriculum(ctx, id); err != nil {
		return err
	}

	// Invalidate cache for this curriculum
	curriculumCacheKey := cache.GenerateCurriculumCacheKey(id.String())
	curriculumBodyCacheKey := cache.GenerateCurriculumBodyCacheKey(id.String())

	// Delete curriculum cache
	if err := cu.cacheService.Delete(ctx, curriculumCacheKey); err != nil {
		cu.logger.Warn("Failed to invalidate curriculum cache after deletion",
			zap.Error(err),
			zap.String("curriculum_id", id.String()))
	} else {
		cu.logger.Debug("Curriculum cache invalidated after deletion", zap.String("curriculum_id", id.String()))
	}

	// Delete curriculum body cache
	if err := cu.cacheService.Delete(ctx, curriculumBodyCacheKey); err != nil {
		cu.logger.Warn("Failed to invalidate curriculum body cache after deletion",
			zap.Error(err),
			zap.String("curriculum_id", id.String()))
	} else {
		cu.logger.Debug("Curriculum body cache invalidated after deletion", zap.String("curriculum_id", id.String()))
	}

	return nil
}
