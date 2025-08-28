package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

// ValidationError represents a formatted validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents the validation error response
type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

// ValidatePhoneNumber validates a phone number using the nyaruka/phonenumbers library
func ValidatePhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}

	// Try to parse the phone number with different country codes
	// We'll try common country codes if no country code is provided
	countryCodes := []string{"BR", "US", "PT", "ES", "FR", "DE", "IT", "GB", "CA", "AU"}

	for _, countryCode := range countryCodes {
		num, err := phonenumbers.Parse(phone, countryCode)
		if err == nil && phonenumbers.IsValidNumber(num) {
			return true
		}
	}

	// If no country code is provided, try to parse as international format
	if strings.HasPrefix(phone, "+") {
		num, err := phonenumbers.Parse(phone, "")
		if err == nil && phonenumbers.IsValidNumber(num) {
			return true
		}
	}

	return false
}

// RegisterCustomValidators registers custom validation functions
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("phone", validatePhone)
}

// validatePhone is a custom validator function for phone numbers
func validatePhone(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return ValidatePhoneNumber(phone)
}

// formatValidationError converts validation errors to friendly messages
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

// HandleValidationError handles validation errors and returns a formatted response
func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		formattedError := formatValidationError(validationErrors)
		c.JSON(http.StatusBadRequest, formattedError)
		return
	}

	// If it's not a validation error, return the original error
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
