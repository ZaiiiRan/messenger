package status

type StatusValidationError struct {
	message string
}

func NewStatusValidationError(message string) *StatusValidationError {
	return &StatusValidationError{
		message: message,
	}
}

func (e *StatusValidationError) Error() string {
	return e.message
}
