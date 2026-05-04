package passwordservice

type PasswordServiceError struct {
	message string
}

func newPasswordServiceError(message string) *PasswordServiceError {
	return &PasswordServiceError{message: message}
}

func (e *PasswordServiceError) Error() string {
	return e.message
}

var (
	ErrPasswordNotFound = newPasswordServiceError("service.password.error.password_not_found")
)
