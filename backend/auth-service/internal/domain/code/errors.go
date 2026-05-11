package code

var (
	ErrTooManyConfirmationAttempts    = NewCodeValidationError("domain.code.validation.error.too_many_confirmation_attempts")
	ErrCodeResendExhausted            = NewCodeValidationError("domain.code.validation.error.number_of_code_resends_has_been_exhausted")
	ErrWaitBeforeNewCode              = NewCodeValidationError("domain.code.validation.error.wait_before_requesting_new_code")
	ErrCodeExpired                    = NewCodeValidationError("domain.code.validation.error.code_has_been_expired")
	ErrTooManyFailedAttempts          = NewCodeValidationError("domain.code.validation.error.too_many_failed_confirmation_attempts")
	ErrLinkExpired                    = NewCodeValidationError("domain.code.validation.error.link_has_expired")
	ErrInvalidCode                    = NewCodeValidationError("domain.code.validation.error.invalid_code")
	ErrInvalidToken                   = NewCodeValidationError("domain.code.validation.error.invalid_token")
	ErrInvalidOrExpiredActivationLink = NewCodeValidationError("domain.code.validation.error.invalid_or_expired_activation_link")
	ErrInvalidOrExpiredResetLink      = NewCodeValidationError("domain.code.validation.error.invalid_or_expired_reset_link")
)
