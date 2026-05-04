package profile

import (
	"regexp"
	"time"
	"unicode/utf8"
)

func validatePhone(phone string) error {
	if phone == "" {
		return nil
	}

	phoneRegex := regexp.MustCompile(`^\+\d{1,3}[-\s]?\(?\d{1,4}\)?(?:[-\s]?\d{1,4}){1,4}$`)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhoneFormat
	}

	phoneDigitRegex := regexp.MustCompile(`\D`)

	digits := phoneDigitRegex.ReplaceAllString(phone, "")
	if len(digits) < 8 || len(digits) > 15 {
		return ErrInvalidPhoneLength
	}

	return nil
}

func validateName(name string, prefix string) error {
	if name == "" {
		return getNameIsEmptyError(prefix)
	}
	if utf8.RuneCountInString(name) < 2 {
		return getNameTooShortError(prefix)
	}

	if utf8.RuneCountInString(name) > 50 {
		return getNameTooLongError(prefix)
	}

	nameRegex := regexp.MustCompile(`^\p{Lu}[\p{L}'’]+(?:[-\s’'][\p{Lu}\p{Ll}][\p{L}'’]*)*$`)
	if !nameRegex.MatchString(name) {
		return getInvalidNameFormatError(prefix)
	}
	return nil
}

func validateBirthdate(birthdate time.Time) error {
	if birthdate.After(time.Now()) {
		return ErrBirthdateInFuture
	}
	return nil
}

func validateBio(bio string) error {
	if bio == "" {
		return nil
	}

	if utf8.RuneCountInString(bio) > 1000 {
		return ErrBioTooLong
	}
	return nil
}
