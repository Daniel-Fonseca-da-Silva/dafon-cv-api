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
	GetUsersWithPagination(ctx context.Context, cursor *uuid.UUID, limit int) ([]dto.UserResponse, dto.CursorPagination, error)
	GetUserDetail(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	ToggleAdmin(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetCurriculumsWithPagination(ctx context.Context, cursor *uuid.UUID, limit int) ([]dto.CurriculumResponse, dto.CursorPagination, error)
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

// GetUsersWithPagination returns cursor-paginated users.
func (uc *adminUseCase) GetUsersWithPagination(ctx context.Context, cursor *uuid.UUID, limit int) ([]dto.UserResponse, dto.CursorPagination, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, hasNextPage, err := uc.userRepo.GetPageAfterID(ctx, cursor, limit)
	if err != nil {
		return nil, dto.CursorPagination{}, fmt.Errorf("get users page: %w", err)
	}

	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		responses[i] = userModelToResponse(u)
	}

	pagination := dto.CursorPagination{
		Limit:       limit,
		HasNextPage: hasNextPage,
	}
	if cursor != nil && *cursor != uuid.Nil {
		cursorStr := cursor.String()
		pagination.Cursor = &cursorStr
	}
	if hasNextPage && len(users) > 0 {
		nextCursor := users[len(users)-1].ID.String()
		pagination.NextCursor = &nextCursor
	}

	return responses, pagination, nil
}

// GetUserDetail returns a single user by ID
func (uc *adminUseCase) GetUserDetail(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user detail %s: %w", id.String(), err)
	}
	resp := userModelToResponse(*user)
	return &resp, nil
}

// ToggleAdmin flips the admin flag for a user
func (uc *adminUseCase) ToggleAdmin(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.ToggleAdmin(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("toggle admin %s: %w", id.String(), err)
	}
	resp := userModelToResponse(*user)
	return &resp, nil
}

// GetCurriculumsWithPagination returns cursor-paginated curriculums.
func (uc *adminUseCase) GetCurriculumsWithPagination(ctx context.Context, cursor *uuid.UUID, limit int) ([]dto.CurriculumResponse, dto.CursorPagination, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	curriculums, hasNextPage, err := uc.curriculumRepo.GetPageAfterID(ctx, cursor, limit)
	if err != nil {
		return nil, dto.CursorPagination{}, fmt.Errorf("get curriculums page: %w", err)
	}

	responses := make([]dto.CurriculumResponse, len(curriculums))
	for i, c := range curriculums {
		responses[i] = curriculumModelToResponse(c)
	}

	pagination := dto.CursorPagination{
		Limit:       limit,
		HasNextPage: hasNextPage,
	}
	if cursor != nil && *cursor != uuid.Nil {
		cursorStr := cursor.String()
		pagination.Cursor = &cursorStr
	}
	if hasNextPage && len(curriculums) > 0 {
		nextCursor := curriculums[len(curriculums)-1].ID.String()
		pagination.NextCursor = &nextCursor
	}

	return responses, pagination, nil
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
