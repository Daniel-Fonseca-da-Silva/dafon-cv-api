package validation

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// IsValidPhone validates a phone number using the nyaruka/phonenumbers library
func IsValidPhone(phone string) bool {
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
