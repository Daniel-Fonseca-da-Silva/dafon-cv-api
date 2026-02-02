package handlers

import (
	"errors"
	"net/http"
	"strconv"

	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminHandler handles HTTP requests for admin (back office) operations
type AdminHandler struct {
	adminUseCase usecases.AdminUseCase
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(adminUseCase usecases.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUseCase}
}

// GetDashboard handles GET /api/v1/admin/dashboard
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	dashboard, err := h.adminUseCase.GetDashboard(c.Request.Context())
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

// GetUsers handles GET /api/v1/admin/users?page=1&page_size=2
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	users, total, err := h.adminUseCase.GetUsersWithPagination(c.Request.Context(), page, pageSize)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetUserDetail handles GET /api/v1/admin/users/:id/detail
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	user, err := h.adminUseCase.GetUserDetail(c.Request.Context(), id)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("user not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// ToggleAdmin handles PATCH /api/v1/admin/users/:id/toggle-admin
func (h *AdminHandler) ToggleAdmin(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	user, err := h.adminUseCase.ToggleAdmin(c.Request.Context(), id)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurriculums handles GET /api/v1/admin/curriculums?page=1&page_size=10
func (h *AdminHandler) GetCurriculums(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "DESC")
	if sortOrder != "ASC" && sortOrder != "DESC" {
		transporthttp.HandleValidationError(c, errors.New("invalid sort_order, must be ASC or DESC"))
		return
	}

	curriculums, total, err := h.adminUseCase.GetCurriculumsWithPagination(c.Request.Context(), page, pageSize, sortBy, sortOrder)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": curriculums,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

// GetCurriculumsStats handles GET /api/v1/admin/curriculums/stats
func (h *AdminHandler) GetCurriculumsStats(c *gin.Context) {
	stats, err := h.adminUseCase.GetCurriculumsStats(c.Request.Context())
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetUsersStats handles GET /api/v1/admin/users/stats
func (h *AdminHandler) GetUsersStats(c *gin.Context) {
	stats, err := h.adminUseCase.GetUsersStats(c.Request.Context())
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

func parsePagination(c *gin.Context) (page, pageSize int, err error) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err = strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid page, must be a positive integer")
	}

	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return 0, 0, errors.New("invalid page_size, must be a positive integer")
	}

	return page, pageSize, nil
}
