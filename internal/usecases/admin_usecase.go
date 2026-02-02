package usecases

import (
	"context"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AdminUseCase defines the interface for admin (back office) operations
type AdminUseCase interface {
	GetDashboard(ctx context.Context) (*dto.DashboardResponse, error)
	GetUsersWithPagination(ctx context.Context, page, pageSize int) ([]dto.UserResponse, int64, error)
	GetUserDetail(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	ToggleAdmin(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetCurriculumsWithPagination(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]dto.CurriculumResponse, int64, error)
	GetCurriculumsStats(ctx context.Context) (*dto.CurriculumsStatsResponse, error)
	GetUsersStats(ctx context.Context) (*dto.UsersStatsResponse, error)
}

type adminUseCase struct {
	userRepo       repositories.UserRepository
	curriculumRepo repositories.CurriculumRepository
	logger         *zap.Logger
}

// NewAdminUseCase creates a new AdminUseCase
func NewAdminUseCase(userRepo repositories.UserRepository, curriculumRepo repositories.CurriculumRepository, logger *zap.Logger) AdminUseCase {
	return &adminUseCase{
		userRepo:       userRepo,
		curriculumRepo: curriculumRepo,
		logger:         logger,
	}
}

// GetDashboard returns dashboard summary (users and curriculums count)
func (uc *adminUseCase) GetDashboard(ctx context.Context) (*dto.DashboardResponse, error) {
	usersCount, err := uc.userRepo.Count(ctx)
	if err != nil {
		uc.logger.Error("Failed to count users for dashboard", zap.Error(err))
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	curriculumsCount, err := uc.curriculumRepo.Count(ctx)
	if err != nil {
		uc.logger.Error("Failed to count curriculums for dashboard", zap.Error(err))
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	return &dto.DashboardResponse{
		UsersCount:       usersCount,
		CurriculumsCount: curriculumsCount,
	}, nil
}

// GetUsersWithPagination returns paginated users and total count
func (uc *adminUseCase) GetUsersWithPagination(ctx context.Context, page, pageSize int) ([]dto.UserResponse, int64, error) {
	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	users, err := uc.userRepo.GetWithPagination(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		responses[i] = userModelToResponse(u)
	}
	return responses, total, nil
}

// GetUserDetail returns a single user by ID
func (uc *adminUseCase) GetUserDetail(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := userModelToResponse(*user)
	return &resp, nil
}

// ToggleAdmin flips the admin flag for a user
func (uc *adminUseCase) ToggleAdmin(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.ToggleAdmin(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := userModelToResponse(*user)
	return &resp, nil
}

// GetCurriculumsWithPagination returns paginated curriculums and total count
func (uc *adminUseCase) GetCurriculumsWithPagination(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]dto.CurriculumResponse, int64, error) {
	total, err := uc.curriculumRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count curriculums: %w", err)
	}

	curriculums, err := uc.curriculumRepo.GetAll(ctx, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.CurriculumResponse, len(curriculums))
	for i, c := range curriculums {
		responses[i] = curriculumModelToResponse(c)
	}
	return responses, total, nil
}

// GetCurriculumsStats returns curriculum statistics
func (uc *adminUseCase) GetCurriculumsStats(ctx context.Context) (*dto.CurriculumsStatsResponse, error) {
	count, err := uc.curriculumRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get curriculums stats: %w", err)
	}
	return &dto.CurriculumsStatsResponse{Total: count}, nil
}

// GetUsersStats returns user statistics
func (uc *adminUseCase) GetUsersStats(ctx context.Context) (*dto.UsersStatsResponse, error) {
	count, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users stats: %w", err)
	}
	return &dto.UsersStatsResponse{Total: count}, nil
}

func userModelToResponse(u models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		ImageURL:   u.ImageURL,
		Country:    u.Country,
		State:      u.State,
		City:       u.City,
		Phone:      u.Phone,
		Employment: u.Employment,
		Gender:     u.Gender,
		Age:        u.Age,
		Salary:     u.Salary,
		Migration:  u.Migration,
		Admin:      u.Admin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func curriculumModelToResponse(c models.Curriculums) dto.CurriculumResponse {
	works := make([]dto.WorkResponse, len(c.Works))
	for i, w := range c.Works {
		works[i] = dto.WorkResponse{
			ID:          w.ID,
			Position:    w.Position,
			Company:     w.Company,
			Description: w.Description,
			StartDate:   w.StartDate,
			EndDate:     w.EndDate,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		}
	}
	educations := make([]dto.EducationResponse, len(c.Educations))
	for i, e := range c.Educations {
		educations[i] = dto.EducationResponse{
			ID:          e.ID,
			Institution: e.Institution,
			Degree:      e.Degree,
			StartDate:   e.StartDate,
			EndDate:     e.EndDate,
			Description: e.Description,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		}
	}
	return dto.CurriculumResponse{
		ID:            c.ID,
		FullName:      c.FullName,
		Email:         c.Email,
		Phone:         c.Phone,
		DriverLicense: c.DriverLicense,
		Intro:         c.Intro,
		Skills:        c.Skills,
		Languages:     c.Languages,
		Courses:       c.Courses,
		SocialLinks:   c.SocialLinks,
		ImageURL:      c.ImageURL,
		Works:         works,
		Educations:    educations,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}
