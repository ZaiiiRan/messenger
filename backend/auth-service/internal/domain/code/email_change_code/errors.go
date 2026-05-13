package emailchangecode

import "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"

var (
	ErrEmptyEmail         = code.NewCodeValidationError("domain.email_change_code.validation.error.email_is_empty")
	ErrInvalidEmailFormat = code.NewCodeValidationError("domain.email_change_code.validation.error.invalid_email_format")
	ErrEmailTooLong       = code.NewCodeValidationError("domain.email_change_code.validation.error.email_too_long")
	ErrSameEmail          = code.NewCodeValidationError("domain.email_change_code.validation.error.same_email")
)
