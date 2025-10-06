package usecases

import (
	"context"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CurriculumUsecase defines the interface for curriculum business logic operations
type CurriculumUseCase interface {
	CreateCurriculum(ctx context.Context, userID uuid.UUID, req *dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error)
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
		UserID:        userID,
	}

	// Validate the curriculum model
	validate := validator.New()
	if err := validate.Struct(curriculum); err != nil {
		return nil, err
	}

	// Create works associated with curriculum
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

	// Create educations associated with curriculum
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

	// Save to database (GORM will handle the foreign key relationship)
	if err := cu.curriculumRepo.Create(ctx, curriculum); err != nil {
		return nil, err
	}

	// Prepare works response
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

	// Prepare educations response
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
		Works:         worksResponse,
		Educations:    educationsResponse,
		CreatedAt:     curriculum.CreatedAt,
		UpdatedAt:     curriculum.UpdatedAt,
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
		Works:         worksResponse,
		Educations:    educationsResponse,
		CreatedAt:     curriculum.CreatedAt,
		UpdatedAt:     curriculum.UpdatedAt,
	}, nil
}
