package password

import (
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewPasswordValidationError("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

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

	if !hasUpper {
		return NewPasswordValidationError("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return NewPasswordValidationError("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return NewPasswordValidationError("password must contain at least one digit")
	}
	if !hasSpecial {
		return NewPasswordValidationError("password must contain at least one special character")
	}

	return nil
}
