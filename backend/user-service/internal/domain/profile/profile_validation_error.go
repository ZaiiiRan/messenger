package profile

type ProfileValidationError struct {
	message string
}

func NewProfileValidationError(message string) *ProfileValidationError {
	return &ProfileValidationError{
		message: message,
	}
}

func (e *ProfileValidationError) Error() string {
	return e.message
}
