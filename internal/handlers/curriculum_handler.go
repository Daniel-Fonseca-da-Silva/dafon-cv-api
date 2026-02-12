package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CurriculumHandler handles HTTP requests for curriculum operations
type CurriculumHandler struct {
	curriculumUseCase usecases.CurriculumUseCase
	userUseCase       usecases.UserUseCase
	logger            *zap.Logger
}

// NewCurriculumHandler creates a new instance of CurriculumHandler
func NewCurriculumHandler(curriculumUseCase usecases.CurriculumUseCase, userUseCase usecases.UserUseCase, logger *zap.Logger) *CurriculumHandler {
	return &CurriculumHandler{
		curriculumUseCase: curriculumUseCase,
		userUseCase:       userUseCase,
		logger:            logger,
	}
}

// CreateCurriculum godoc
// @Summary      Create curriculum
// @Description  Creates a new curriculum for a user
// @Tags         curriculum
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateCurriculumRequest  true  "Create payload"
// @Success      201   {object}  dto.CurriculumResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error or invalid user ID"
// @Failure      404   {object}  dto.ErrorResponse  "User not found"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/curriculums [post]
// @Security     BearerAuth
func (h *CurriculumHandler) CreateCurriculum(c *gin.Context) {
	var req dto.CreateCurriculumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	// Validate phone number using our custom validation
	if !validation.IsValidPhone(req.Phone) {
		c.JSON(http.StatusBadRequest, transporthttp.ValidationErrorResponse{
			Message: "Invalid input data",
			Errors: []transporthttp.ValidationError{
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
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	// Verify if the user exists in the database
	_, err = h.userUseCase.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		h.abortWithInternalServerError(c, "verify user for curriculum creation", err)
		return
	}

	curriculum, err := h.curriculumUseCase.CreateCurriculum(c.Request.Context(), userUUID, &req)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			transporthttp.HandleValidationError(c, err)
			return
		}
		h.abortWithInternalServerError(c, "create curriculum", err)
		return
	}

	c.JSON(http.StatusCreated, curriculum)
}

// GetCurriculumByID godoc
// @Summary      Get curriculum by ID
// @Description  Returns a curriculum by ID
// @Tags         curriculum
// @Accept       json
// @Produce      json
// @Param        curriculum_id  path      string  true  "Curriculum ID"
// @Success      200            {object}  dto.CurriculumResponse
// @Failure      400            {object}  dto.ErrorResponseValidation  "Invalid curriculum ID format"
// @Failure      404            {object}  dto.ErrorResponse  "Curriculum not found"
// @Failure      500            {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/curriculums/{curriculum_id} [get]
// @Security     BearerAuth
func (h *CurriculumHandler) GetCurriculumByID(c *gin.Context) {
	idStr := c.Param("curriculum_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	curriculum, err := h.curriculumUseCase.GetCurriculumByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "curriculum not found")
			return
		}
		h.abortWithInternalServerError(c, "get curriculum by id", err)
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

// GetAllCurriculums godoc
// @Summary      Get all curriculums by user ID
// @Description  Returns paginated list of curriculums for the given user
// @Tags         curriculum
// @Accept       json
// @Produce      json
// @Param        user_id   path      string  true   "User ID"
// @Param        cursor    query     string  false  "Cursor (UUID) to fetch items after"
// @Param        limit     query     int     false  "Items per page" default(10)
// @Success      200       {object}  dto.CurriculumListResponse
// @Failure      400       {object}  dto.ErrorResponseValidation  "Invalid user ID or pagination params"
// @Failure      404       {object}  dto.ErrorResponse  "User not found"
// @Failure      500       {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/curriculums/get-all-by-user/{user_id} [get]
// @Security     BearerAuth
func (h *CurriculumHandler) GetAllCurriculums(c *gin.Context) {
	// Extrair user_id dos parâmetros da URL
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	// Verificar se o usuário existe
	_, err = h.userUseCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		h.abortWithInternalServerError(c, "verify user for curriculum list", err)
		return
	}

	cursor, limit, err := parseCursorPagination(c)
	if err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	curriculums, pagination, err := h.curriculumUseCase.GetAllCurriculums(c.Request.Context(), userID, cursor, limit)
	if err != nil {
		h.abortWithInternalServerError(c, "get all curriculums", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       curriculums,
		"pagination": pagination,
	})
}

// GetCurriculumBody godoc
// @Summary      Get curriculum body
// @Description  Returns the curriculum content in text format
// @Tags         curriculum
// @Accept       json
// @Produce      json
// @Param        curriculum_id  path      string  true  "Curriculum ID"
// @Success      200            {object}  dto.CurriculumBodyResponse
// @Failure      400            {object}  dto.ErrorResponseValidation  "Invalid curriculum ID format"
// @Failure      404            {object}  dto.ErrorResponse  "Curriculum not found"
// @Failure      500            {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/curriculums/get-body/{curriculum_id} [get]
// @Security     BearerAuth
func (h *CurriculumHandler) GetCurriculumBody(c *gin.Context) {
	curriculumIDStr := c.Param("curriculum_id")
	curriculumID, err := uuid.Parse(curriculumIDStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	curriculumBody, err := h.curriculumUseCase.GetCurriculumBody(c.Request.Context(), curriculumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "curriculum not found")
			return
		}
		h.abortWithInternalServerError(c, "get curriculum body", err)
		return
	}

	c.JSON(http.StatusOK, curriculumBody)
}

// DeleteCurriculum godoc
// @Summary      Delete curriculum by ID
// @Description  Deletes a curriculum by ID
// @Tags         curriculum
// @Accept       json
// @Produce      json
// @Param        curriculum_id  path      string  true  "Curriculum ID"
// @Success      200            {object}  dto.MessageResponse
// @Failure      400            {object}  dto.ErrorResponseValidation  "Invalid curriculum ID format"
// @Failure      404            {object}  dto.ErrorResponse  "Curriculum not found"
// @Failure      500            {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/curriculums/{curriculum_id} [delete]
// @Security     BearerAuth
func (h *CurriculumHandler) DeleteCurriculum(c *gin.Context) {
	idStr := c.Param("curriculum_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	if err := h.curriculumUseCase.DeleteCurriculum(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "curriculum not found")
			return
		}
		h.abortWithInternalServerError(c, "delete curriculum", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curriculum deleted successfully"})
}

func (h *CurriculumHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("Curriculum handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
