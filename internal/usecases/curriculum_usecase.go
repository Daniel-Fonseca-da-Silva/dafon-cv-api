package usecases

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
)

// CurriculumUsecase defines the interface for curriculum business logic operations
type CurriculumUseCase interface {
	CreateCurriculum(ctx context.Context, req *dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error)
	GetCurriculumByID(ctx context.Context, id uuid.UUID) (*dto.CurriculumResponse, error)
}

// curriculumUsecase implements CurriculumUseCase interface
type curriculumUseCase struct {
	curriculumRepo repositories.CurriculumRepository
}

// NewCurriculumUsecase creates a new instance of CurriculumUsecase
func NewCurriculumUseCase(curriculumRepo repositories.CurriculumRepository) CurriculumUseCase {
	return &curriculumUseCase{
		curriculumRepo: curriculumRepo,
	}
}

// CreateCurriculum creates a new curriculum in the database
func (cu *curriculumUseCase) CreateCurriculum(ctx context.Context, req *dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error) {
	// Create curriculum model
	curriculum := &models.Curriculums{
		FullName:          req.FullName,
		Email:             req.Email,
		DriverLicense:     req.DriverLicense,
		Intro:             req.Intro,
		DateDisponibility: *req.DateDisponibility,
		Languages:         req.Languages,
		LevelEducation:    req.LevelEducation,
		JobDescription:    req.JobDescription,
	}

	// Create works associated with curriculum
	for _, workReq := range req.Works {
		work := models.Work{
			JobTitle:           workReq.JobTitle,
			CompanyName:        workReq.CompanyName,
			CompanyDescription: workReq.CompanyDescription,
			StartDate:          workReq.StartDate,
			EndDate:            workReq.EndDate,
		}
		curriculum.Works = append(curriculum.Works, work)
	}

	// Save to database (GORM will handle the foreign key relationship)
	if err := cu.curriculumRepo.Create(ctx, curriculum); err != nil {
		return nil, err
	}

	// Prepare works response
	worksResponse := make([]dto.WorkResponse, 0, len(curriculum.Works))
	for _, work := range curriculum.Works {
		worksResponse = append(worksResponse, dto.WorkResponse{
			ID:                 work.ID,
			JobTitle:           work.JobTitle,
			CompanyName:        work.CompanyName,
			CompanyDescription: work.CompanyDescription,
			StartDate:          work.StartDate,
			EndDate:            work.EndDate,
			CreatedAt:          work.CreatedAt,
			UpdatedAt:          work.UpdatedAt,
		})
	}

	// Return response
	return &dto.CurriculumResponse{
		ID:                curriculum.ID,
		FullName:          curriculum.FullName,
		Email:             curriculum.Email,
		DriverLicense:     curriculum.DriverLicense,
		Intro:             curriculum.Intro,
		DateDisponibility: curriculum.DateDisponibility,
		Languages:         curriculum.Languages,
		LevelEducation:    curriculum.LevelEducation,
		JobDescription:    curriculum.JobDescription,
		Works:             worksResponse,
		CreatedAt:         curriculum.CreatedAt,
		UpdatedAt:         curriculum.UpdatedAt,
	}, nil
}

// GetCurriculumByID retrieves a curriculum by ID
func (cu *curriculumUseCase) GetCurriculumByID(ctx context.Context, id uuid.UUID) (*dto.CurriculumResponse, error) {
	curriculum, err := cu.curriculumRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert works from models to DTOs
	worksResponse := make([]dto.WorkResponse, 0, len(curriculum.Works))
	for _, work := range curriculum.Works {
		worksResponse = append(worksResponse, dto.WorkResponse{
			ID:                 work.ID,
			JobTitle:           work.JobTitle,
			CompanyName:        work.CompanyName,
			CompanyDescription: work.CompanyDescription,
			StartDate:          work.StartDate,
			EndDate:            work.EndDate,
			CreatedAt:          work.CreatedAt,
			UpdatedAt:          work.UpdatedAt,
		})
	}

	return &dto.CurriculumResponse{
		ID:                curriculum.ID,
		FullName:          curriculum.FullName,
		Email:             curriculum.Email,
		DriverLicense:     curriculum.DriverLicense,
		Intro:             curriculum.Intro,
		DateDisponibility: curriculum.DateDisponibility,
		Languages:         curriculum.Languages,
		LevelEducation:    curriculum.LevelEducation,
		JobDescription:    curriculum.JobDescription,
		Works:             worksResponse,
		CreatedAt:         curriculum.CreatedAt,
		UpdatedAt:         curriculum.UpdatedAt,
	}, nil
}
