package code

type CodeValidationError struct {
	message string
}

func NewCodeValidationError(message string) *CodeValidationError {
	return &CodeValidationError{
		message: message,
	}
}

func (e *CodeValidationError) Error() string {
	return e.message
}
