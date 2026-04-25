package password

type PasswordValidationError struct {
	message string
}

func NewPasswordValidationError(message string) *PasswordValidationError {
	return &PasswordValidationError{
		message: message,
	}
}

func (e *PasswordValidationError) Error() string {
	return e.message
}
