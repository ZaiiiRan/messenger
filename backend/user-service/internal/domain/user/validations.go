package user

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}

	if strings.Contains(username, " ") {
		return fmt.Errorf("username cannot contain spaces")
	}

	if utf8.RuneCountInString(username) < 5 {
		return fmt.Errorf("username must be at least 5 characters")
	}
	return nil
}
