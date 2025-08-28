package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
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
		utils.HandleValidationError(c, err)
		return
	}

	// Validate phone number using our custom validation
	if !utils.ValidatePhoneNumber(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"errors": []utils.ValidationError{
				{
					Field:   "phone",
					Message: "Invalid phone number",
				},
			},
		})
		return
	}

	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.HandleValidationError(c, errors.New("user not authenticated"))
		return
	}
	userUUID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid user ID in context"))
		return
	}

	// Verify if the authenticated user exists in the database
	_, err = h.userUseCase.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		utils.HandleValidationError(c, errors.New("user not found"))
		return
	}

	curriculum, err := h.curriculumUseCase.CreateCurriculum(c.Request.Context(), userUUID, &req)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, curriculum)
}

// GetCurriculumByID handles GET /curriculums/:id request
func (h *CurriculumHandler) GetCurriculumByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	curriculum, err := h.curriculumUseCase.GetCurriculumByID(c.Request.Context(), id)
	if err != nil {
		utils.HandleValidationError(c, errors.New("curriculum not found"))
		return
	}

	c.JSON(http.StatusOK, curriculum)
}
