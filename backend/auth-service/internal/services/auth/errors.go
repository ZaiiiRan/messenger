package authservice

type AuthServiceError struct {
	message string
}

func newAuthServiceError(message string) *AuthServiceError {
	return &AuthServiceError{message: message}
}

func (e *AuthServiceError) Error() string {
	return e.message
}

var (
	ErrInvalidLoginOrPassword = newAuthServiceError("service.auth.error.invalid_login_or_password")
	ErrInvalidCredentials     = newAuthServiceError("service.auth.error.invalid_credentials")
	ErrTooManySessions        = newAuthServiceError("service.auth.error.too_many_sessions")
)
