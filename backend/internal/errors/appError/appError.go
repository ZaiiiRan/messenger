package appError

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	StatusCode int
	Err        error
}

// Get error message
func (e *AppError) Error() string {
	return e.Err.Error()
}

// New AppError object
func NewAppError(statusCode int, errMessage string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Err:        errors.New(errMessage),
	}
}

// Converting error to AppError
func WrapError(statusCode int, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Err:        err,
	}
}

// Bad request error
func BadRequest(errMessage string) *AppError {
	return NewAppError(fiber.StatusBadRequest, errMessage)
}

// Unauthorized error
func Unauthorized(errMessage string) *AppError {
	return NewAppError(fiber.StatusUnauthorized, errMessage)
}

// Not found error
func NotFound(errMessage string) *AppError {
	return NewAppError(fiber.StatusNotFound, errMessage)
}

// Internal server error
func InternalServerError(errMessage string) *AppError {
	return NewAppError(fiber.StatusInternalServerError, errMessage)
}
