package codeservice

type CodeServiceError struct {
	message string
}

func newCodeServiceError(message string) *CodeServiceError {
	return &CodeServiceError{message: message}
}

func (e *CodeServiceError) Error() string {
	return e.message
}

var (
	ErrCodeNotFound = newCodeServiceError("service.code.error.code_not_found")
)
