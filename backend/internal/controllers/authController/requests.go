package authController

import (
	"strings"
)

// Register request format
type RegisterRequest struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Phone     *string `json:"phone"`
	Birthdate *string `json:"birthdate"`
}

// trim spaces in register request
func (r *RegisterRequest) trimSpaces() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.Firstname = strings.TrimSpace(r.Firstname)
	r.Lastname = strings.TrimSpace(r.Lastname)
	if r.Phone != nil {
		trimmedPhone := strings.TrimSpace(*r.Phone)
		r.Phone = &trimmedPhone
	}
	if r.Birthdate != nil {
		trimmedBirthdate := strings.TrimSpace(*r.Birthdate)
		r.Birthdate = &trimmedBirthdate
	}
}

// Login request format
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// trim spaces in login request
func (r *LoginRequest) trimSpaces() {
	r.Login = strings.TrimSpace(r.Login)
	r.Password = strings.TrimSpace(r.Password)
}

// Activate Register request format
type ActivateRequest struct {
	Code string `json:"code"`
}

// trim spaces in activation request
func (r *ActivateRequest) trimSpaces() {
	r.Code = strings.TrimSpace(r.Code)
}
