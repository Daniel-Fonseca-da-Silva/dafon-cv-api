package transporthttp

import (
	"net/http"
	"strings"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		formattedError := formatValidationError(validationErrors)
		c.JSON(http.StatusBadRequest, formattedError)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func ValidatePasswordAndReturnError(c *gin.Context, password string) bool {
	if !validation.IsStrongPassword(password) {
		errorResponse := ValidationErrorResponse{
			Message: "Invalid input data",
			Errors: []ValidationError{
				{
					Field:   "password",
					Message: "Password must be at least 8 characters long and contain uppercase, lowercase, digit, and special character",
				},
			},
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return false
	}
	return true
}

func formatValidationError(err error) ValidationErrorResponse {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			var message string

			switch e.Tag() {
			case "required":
				message = "Field is required"
			case "email":
				message = "Invalid email"
			case "min":
				message = "Value is too short"
			case "max":
				message = "Value is too long"
			case "url":
				message = "Invalid URL"
			case "uuid":
				message = "Invalid ID"
			case "phone":
				message = "Invalid phone number"
			default:
				message = "Invalid value"
			}

			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}

	return ValidationErrorResponse{
		Message: "Invalid input data",
		Errors:  errors,
	}
}
