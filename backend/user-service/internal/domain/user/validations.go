package user

import (
	"regexp"
	"strings"
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

func validateUsername(username string) error {
	if username == "" {
		return ErrEmptyUsername
	}

	if strings.Contains(username, " ") {
		return ErrUsernameContainsSpaces
	}

	if utf8.RuneCountInString(username) < 5 {
		return ErrUsernameTooShort
	}

	if utf8.RuneCountInString(username) > 30 {
		return ErrUsernameTooLong
	}

	return nil
}
