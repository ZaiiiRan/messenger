package profile

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"
)

func validatePhone(phone string) error {
	if phone == "" {
		return nil
	}

	phoneRegex := regexp.MustCompile(`^\+7\(9\d{2}\)-\d{3}-\d{2}-\d{2}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("phone must be in format +7(9xx)-xxx-xx-xx or empty")
	}
	return nil
}

func validateName(name string, prefix string) error {
	if name == "" {
		return fmt.Errorf("%sname is empty", prefix)
	}
	if utf8.RuneCountInString(name) < 2 {
		return fmt.Errorf("%sname must be at least 2 characters", prefix)
	}

	nameRegex := regexp.MustCompile(`^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$`)
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("%sname must start with a capital letter", prefix)
	}
	return nil
}

func validateBirthdate(birthdate time.Time) error {
	if birthdate.After(time.Now()) {
		return fmt.Errorf("birthdate cannot be in the future")
	}
	return nil
}
