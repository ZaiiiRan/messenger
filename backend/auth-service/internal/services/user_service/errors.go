package userservice

type UserServiceError struct {
	message string
}

func newUserServiceError(message string) *UserServiceError {
	return &UserServiceError{message: message}
}

func (e *UserServiceError) Error() string {
	return e.message
}

var (
	ErrUserIsDeleted           = newUserServiceError("service.user.error.user_is_deleted")
	ErrUserAlreadyActivated    = newUserServiceError("service.user.error.user_is_already_activated")
	ErrSameEmail               = newUserServiceError("service.user.error.same_email")
	ErrWaitBeforeEmailChanging = newUserServiceError("service.user.error.wait_before_email_changing")
	ErrUserWithEmailExists     = newUserServiceError("service.user.error.user_with_this_email_exists")
)
