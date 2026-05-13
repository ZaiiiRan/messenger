package emailchangecode

import (
	"regexp"
	"unicode/utf8"
)

func validateEmail(email string) error {
	if email == "" {
		return ErrEmptyEmail
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat
	}

	if utf8.RuneCountInString(email) > 254 {
		return ErrEmailTooLong
	}

	return nil
}
