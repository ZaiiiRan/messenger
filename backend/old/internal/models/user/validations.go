package user

import (
	appErr "backend/internal/errors/appError"
	"regexp"
	"time"
)

// validate username
func validateUsername(username string) error {
	if username == "" {
		return appErr.BadRequest("username is empty")
	}

	candidate, err := GetUserByUsername(username)
	if err != nil && err.Error() != "user not found" {
		return err
	}
	if candidate != nil {
		return appErr.BadRequest("user with the same username already exists")
	} else if len(username) < 5 {
		return appErr.BadRequest("username must be at least 5 characters")
	}
	return nil
}

// validate email
func validateEmail(email string) error {
	if email == "" {
		return appErr.BadRequest("email is empty")
	}

	candidate, err := GetUserByEmail(email)
	if err != nil && err.Error() != "user not found" {
		return err
	}
	if candidate != nil {
		return appErr.BadRequest("user with the same email already exists")
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(email) {
		return appErr.BadRequest("invalid email format")
	}
	return nil
}

// validate phone
func validatePhone(phone string) error {
	if phone == "" {
		return nil
	}

	candidate, err := GetUserByPhone(phone)
	if err != nil && err.Error() != "user not found" {
		return err
	}
	if candidate != nil {
		return appErr.BadRequest("user with the same phone number already exists")
	}

	phoneRegex := regexp.MustCompile(`^\+7\(9\d{2}\)-\d{3}-\d{2}-\d{2}$`)
	if !phoneRegex.MatchString(phone) {
		return appErr.BadRequest("phone must be in format +7(9xx)-xxx-xx-xx or empty")
	}
	return nil
}

// validate names (firstname and lastname)
func validateName(name string) error {
	if name == "" {
		return appErr.BadRequest("name is empty")
	}
	if len(name) < 2 {
		return appErr.BadRequest("name must be at least 2 characters")
	}

	nameRegex := regexp.MustCompile(`^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$`)
	if !nameRegex.MatchString(name) {
		return appErr.BadRequest("name must start with a capital letter")
	}
	return nil
}

// validate password
func validatePassword(password string) error {
	if password == "" {
		return appErr.BadRequest("password is empty")
	}

	var (
		hasUpperCase   = regexp.MustCompile(`[A-ZА-ЯЁ]`).MatchString(password)
		hasLowerCase   = regexp.MustCompile(`[a-zа-яё]`).MatchString(password)
		hasNumber      = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecialChar = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	if len(password) < 8 {
		return appErr.BadRequest("password must be at least 8 characters")
	}
	if !hasUpperCase {
		return appErr.BadRequest("password must contain at least one uppercase letter")
	}
	if !hasLowerCase {
		return appErr.BadRequest("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return appErr.BadRequest("password must contain at least one digit")
	}
	if !hasSpecialChar {
		return appErr.BadRequest("password must contain at least one special character")
	}
	return nil
}

// validate birthdate
func validateBirthdate(birthdate *time.Time) error {
	if birthdate.After(time.Now()) {
		return appErr.BadRequest("birthdate cannot be in the future")
	}
	return nil
}
