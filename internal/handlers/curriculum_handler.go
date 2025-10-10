package handlers

import (
	"errors"
	"net/http"
	"strconv"

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

	// Get user ID from request body
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	// Verify if the user exists in the database
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
	idStr := c.Param("curriculum_id")
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

// GetAllCurriculums traz todos os curriculums paginados de um usuário específico
func (h *CurriculumHandler) GetAllCurriculums(c *gin.Context) {
	// Extrair user_id dos parâmetros da URL
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	// Verificar se o usuário existe
	_, err = h.userUseCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.HandleValidationError(c, errors.New("user not found"))
		return
	}

	// Obter parâmetros de query com valores padrão
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "DESC")

	// Converter page para int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.HandleValidationError(c, errors.New("invalid page format, must be a positive integer"))
		return
	}

	// Converter pageSize para int
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		utils.HandleValidationError(c, errors.New("invalid page_size format, must be a positive integer"))
		return
	}

	// Validar sortOrder
	if sortOrder != "ASC" && sortOrder != "DESC" {
		utils.HandleValidationError(c, errors.New("invalid sort_order format, must be ASC or DESC"))
		return
	}

	// Chamar o usecase passando o userID
	curriculums, err := h.curriculumUseCase.GetAllCurriculums(c.Request.Context(), userID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": curriculums,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

// GetCurriculumBody handles GET /curriculums/get-body/:curriculum_id request
func (h *CurriculumHandler) GetCurriculumBody(c *gin.Context) {
	curriculumIDStr := c.Param("curriculum_id")
	curriculumID, err := uuid.Parse(curriculumIDStr)
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	curriculumBody, err := h.curriculumUseCase.GetCurriculumBody(c.Request.Context(), curriculumID)
	if err != nil {
		utils.HandleValidationError(c, errors.New("curriculum not found"))
		return
	}

	c.JSON(http.StatusOK, curriculumBody)
}

// DeleteCurriculum Deleta um curriculum por ID
func (h *CurriculumHandler) DeleteCurriculum(c *gin.Context) {
	idStr := c.Param("curriculum_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	if err := h.curriculumUseCase.DeleteCurriculum(c.Request.Context(), id); err != nil {
		utils.HandleValidationError(c, errors.New("failed to delete curriculum"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curriculum deleted successfully"})
}
