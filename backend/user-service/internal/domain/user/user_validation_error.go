package user

type UserValidationError struct {
	message string
}

func NewUserValidationError(message string) *UserValidationError {
	return &UserValidationError{
		message: message,
	}
}

func (e *UserValidationError) Error() string {
	return e.message
}
