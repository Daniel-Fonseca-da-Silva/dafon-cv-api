package handlers

import (
	"errors"
	"net/http"
	"strconv"

	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AdminHandler handles HTTP requests for admin (back office) operations
type AdminHandler struct {
	adminUseCase usecases.AdminUseCase
	logger       *zap.Logger
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(adminUseCase usecases.AdminUseCase, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		adminUseCase: adminUseCase,
		logger:       logger,
	}
}

// GetDashboard godoc
// @Summary      Get admin dashboard
// @Description  Returns dashboard summary (users and curriculums count). Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.DashboardResponse
// @Failure      400  {object}  dto.ErrorResponseValidation  "Bad request"
// @Failure      401  {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403  {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/dashboard [get]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	dashboard, err := h.adminUseCase.GetDashboard(c.Request.Context())
	if err != nil {
		h.abortWithInternalServerError(c, "get dashboard", err)
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

// GetUsers godoc
// @Summary      Get users (paginated)
// @Description  Returns paginated list of users. Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        cursor    query     string  false  "Cursor (UUID) to fetch items after"
// @Param        limit     query     int     false  "Items per page" default(10)
// @Success      200       {object}  dto.AdminUsersListResponse
// @Failure      400       {object}  dto.ErrorResponseValidation  "Invalid pagination params"
// @Failure      401       {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403       {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500       {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/users [get]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) GetUsers(c *gin.Context) {
	cursor, limit, err := parseCursorPagination(c)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	users, pagination, err := h.adminUseCase.GetUsersWithPagination(c.Request.Context(), cursor, limit)
	if err != nil {
		h.abortWithInternalServerError(c, "get users", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       users,
		"pagination": pagination,
	})
}

// GetUserDetail godoc
// @Summary      Get user detail by ID
// @Description  Returns a single user by ID. Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponseValidation  "Invalid user ID format"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      401  {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403  {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/users/{id}/detail [get]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	user, err := h.adminUseCase.GetUserDetail(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		h.abortWithInternalServerError(c, "get user detail", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// ToggleAdmin godoc
// @Summary      Toggle user admin flag
// @Description  Toggles the admin flag for a user. Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponseValidation  "Invalid user ID format"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      401  {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403  {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/users/{id}/toggle-admin [patch]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) ToggleAdmin(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	user, err := h.adminUseCase.ToggleAdmin(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		h.abortWithInternalServerError(c, "toggle admin", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurriculums godoc
// @Summary      Get curriculums (paginated)
// @Description  Returns paginated list of curriculums. Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        cursor    query     string  false  "Cursor (UUID) to fetch items after"
// @Param        limit     query     int     false  "Items per page" default(10)
// @Success      200        {object}  dto.AdminCurriculumsListResponse
// @Failure      400        {object}  dto.ErrorResponseValidation  "Invalid pagination or sort params"
// @Failure      401        {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403        {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500        {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/curriculums [get]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) GetCurriculums(c *gin.Context) {
	cursor, limit, err := parseCursorPagination(c)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	curriculums, pagination, err := h.adminUseCase.GetCurriculumsWithPagination(c.Request.Context(), cursor, limit)
	if err != nil {
		h.abortWithInternalServerError(c, "get curriculums", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       curriculums,
		"pagination": pagination,
	})
}

// GetCurriculumsStats godoc
// @Summary      Get curriculums statistics
// @Description  Returns curriculums count. Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.CurriculumsStatsResponse
// @Failure      401  {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403  {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/curriculums/stats [get]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) GetCurriculumsStats(c *gin.Context) {
	stats, err := h.adminUseCase.GetCurriculumsStats(c.Request.Context())
	if err != nil {
		h.abortWithInternalServerError(c, "get curriculums stats", err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetUsersStats godoc
// @Summary      Get users statistics
// @Description  Returns users count. Requires admin user (X-User-ID header).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.UsersStatsResponse
// @Failure      401  {object}  dto.ErrorResponse  "X-User-ID header required"
// @Failure      403  {object}  dto.ErrorResponse  "Admin access required"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/admin/users/stats [get]
// @Security     BearerAuth
// @Security     UserIDHeader
func (h *AdminHandler) GetUsersStats(c *gin.Context) {
	stats, err := h.adminUseCase.GetUsersStats(c.Request.Context())
	if err != nil {
		h.abortWithInternalServerError(c, "get users stats", err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

func parseCursorPagination(c *gin.Context) (*uuid.UUID, int, error) {
	const (
		defaultLimit = 10
		maxLimit     = 100
	)

	limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return nil, 0, errors.New("invalid limit, must be a positive integer")
	}
	if limit > maxLimit {
		return nil, 0, errors.New("invalid limit, must be <= 100")
	}

	cursorStr := c.Query("cursor")
	if cursorStr == "" {
		return nil, limit, nil
	}

	cursor, err := uuid.Parse(cursorStr)
	if err != nil {
		return nil, 0, errors.New("invalid cursor, must be a valid UUID")
	}

	return &cursor, limit, nil
}

func (h *AdminHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Admin handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
