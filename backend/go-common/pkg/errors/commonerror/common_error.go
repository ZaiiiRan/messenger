package commonerror

type CommonError struct {
	message string
}

func (e *CommonError) Error() string {
	return e.message
}

func NewCommonError(message string) *CommonError {
	return &CommonError{
		message: message,
	}
}

var (
	ErrCaceled            = NewCommonError("error.canceled")
	ErrUnknown            = NewCommonError("error.unknown")
	ErrInvalidArgument    = NewCommonError("error.invalid_argument")
	ErrDeadlineExceeded   = NewCommonError("error.deadline_exceeded")
	ErrNotFound           = NewCommonError("error.not_found")
	ErrAlreadyExists      = NewCommonError("error.already_exists")
	ErrPermissionDenied   = NewCommonError("error.permission_denied")
	ErrResourceExhausted  = NewCommonError("error.resource_exhausted")
	ErrFailedPrecondition = NewCommonError("error.failed_precondition")
	ErrAborted            = NewCommonError("error.aborted")
	ErrOutOfRange         = NewCommonError("error.out_of_range")
	ErrUnimplemented      = NewCommonError("error.unimplemented")
	ErrInternal           = NewCommonError("error.internal_server_error")
	ErrUnavailable        = NewCommonError("error.unavailable")
	ErrDataLoss           = NewCommonError("error.data_loss")
	ErrUnauthorized       = NewCommonError("error.unauthorized")
)
