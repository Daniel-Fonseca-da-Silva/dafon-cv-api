package utils

import (
	"net/http"
	"strings"
	"unicode"

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

// ValidateStrongPassword validates if a password meets strong password requirements
// Requirements:
// - At least 8 characters long
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
// - At least one special character
func ValidateStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// ValidatePasswordAndReturnError validates a password and returns a formatted error response if invalid
func ValidatePasswordAndReturnError(c *gin.Context, password string) bool {
	if !ValidateStrongPassword(password) {
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
