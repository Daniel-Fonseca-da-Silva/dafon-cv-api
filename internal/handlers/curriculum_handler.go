package handlers

import (
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CurriculumHandler handles HTTP requests for curriculum operations
type CurriculumHandler struct {
	curriculumUseCase usecases.CurriculumUseCase
	userUseCase       usecases.UserUseCase
}

// NewCurriculumHandler creates a new instance of CurriculumHandler
func NewCurriculumHandler(curriculumUseCase usecases.CurriculumUseCase, userUseCase usecases.UserUseCase) *CurriculumHandler {
	return &CurriculumHandler{
		curriculumUseCase: curriculumUseCase,
		userUseCase:       userUseCase,
	}
}

// CreateCurriculum handles POST /curriculums request
func (h *CurriculumHandler) CreateCurriculum(c *gin.Context) {
	var req dto.CreateCurriculumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}

	// Verify if the authenticated user exists in the database
	_, err = h.userUseCase.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	curriculum, err := h.curriculumUseCase.CreateCurriculum(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, curriculum)
}

// GetCurriculumByID handles GET /curriculums/:id request
func (h *CurriculumHandler) GetCurriculumByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid curriculum ID format"})
		return
	}

	curriculum, err := h.curriculumUseCase.GetCurriculumByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curriculum not found"})
	}

	c.JSON(http.StatusOK, curriculum)
}
