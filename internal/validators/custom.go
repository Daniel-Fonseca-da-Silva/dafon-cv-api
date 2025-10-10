package validators

import (
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/validation"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers custom validation functions
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("strong_password", validateStrongPassword)
}

// validatePhone is a custom validator function for phone numbers
func validatePhone(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return validation.IsValidPhone(phone)
}

// validateStrongPassword is a custom validator function for strong passwords
func validateStrongPassword(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return validation.IsStrongPassword(password)
}
