package user

var (
	ErrEmptyEmail             = NewUserValidationError("domain.user.validation.error.email_is_empty")
	ErrInvalidEmailFormat     = NewUserValidationError("domain.user.validation.error.invalid_email_format")
	ErrEmailTooLong           = NewUserValidationError("domain.user.validation.error.email_too_long")
	ErrEmptyUsername          = NewUserValidationError("domain.user.validation.error.username_is_empty")
	ErrUsernameContainsSpaces = NewUserValidationError("domain.user.validation.error.username_contains_spaces")
	ErrUsernameTooShort       = NewUserValidationError("username must be at least 5 characters long")
	ErrUsernameTooLong        = NewUserValidationError("domain.user.validation.error.username_too_long")
)
