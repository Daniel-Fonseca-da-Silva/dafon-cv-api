package handlers

import (
	"errors"
	"net/http"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	apperrors "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	transporthttp "github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/transport/http"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userUseCase usecases.UserUseCase
	logger      *zap.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userUseCase usecases.UserUseCase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Get a user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponseValidation  "Invalid user ID format"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/user/{id} [get]
// @Security     BearerAuth
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	user, err := h.userUseCase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		h.abortWithInternalServerError(c, "get user by id", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Returns a list of all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.UsersResponse
// @Failure      400  {object}  dto.ErrorResponseValidation  "Bad request"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/user/all [get]
// @Security     BearerAuth
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userUseCase.GetAllUsers(c.Request.Context())
	if err != nil {
		h.abortWithInternalServerError(c, "get all users", err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary      Update user by ID
// @Description  Updates an existing user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "User ID"
// @Param        body  body      dto.UpdateUserRequest   true  "Update payload"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Invalid user ID format or validation error"
// @Failure      404   {object}  dto.ErrorResponse  "User not found"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/user/{id} [patch]
// @Security     BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	user, err := h.userUseCase.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			transporthttp.HandleValidationError(c, errors.New("user already exists"))
			return
		}
		h.abortWithInternalServerError(c, "update user", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create user
// @Description  Creates a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserRequest  true  "Create payload"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponseValidation  "Validation error"
// @Failure      500   {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.HandleValidationError(c, err)
		return
	}

	user, err := h.userUseCase.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			transporthttp.HandleValidationError(c, errors.New("user already exists"))
			return
		}
		h.abortWithInternalServerError(c, "create user", err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// DeleteUser godoc
// @Summary      Delete user by ID
// @Description  Deletes a user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponseValidation  "Invalid user ID format"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      500  {object}  dto.ErrorResponseServer  "Internal server error"
// @Router       /api/v1/user/{id} [delete]
// @Security     BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		transporthttp.HandleValidationError(c, errors.New("invalid user ID format"))
		return
	}

	if err := h.userUseCase.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			transporthttp.HandleUseCaseError(c, err, "user not found")
			return
		}
		h.abortWithInternalServerError(c, "delete user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) abortWithInternalServerError(c *gin.Context, operation string, err error) {
	if h.logger != nil {
		h.logger.Error("User handler failed",
			zap.String("operation", operation),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
