package userservice

type UserServiceError struct {
	message string
}

func newUserServiceError(message string) *UserServiceError {
	return &UserServiceError{
		message: message,
	}
}

func (e *UserServiceError) Error() string {
	return e.message
}

var (
	ErrUserWithUsernameExists       = newUserServiceError("service.user.error.user_with_this_username_exists")
	ErrUserWithEmailExists          = newUserServiceError("service.user.error.user_with_this_email_exists")
	ErrUserIdIsRequired             = newUserServiceError("service.user.error.user_id_is_required")
	ErrUsernameIsRequired           = newUserServiceError("service.user.error.username_is_required")
	ErrEmailIsRequired              = newUserServiceError("service.user.error.email_is_required")
	ErrUserNotFound                 = newUserServiceError("service.user.error.user_not_found")
	ErrUserIsAlreadyActivated       = newUserServiceError("service.user.error.user_is_already_activated")
	ErrDeletedUserCannotBeActivated = newUserServiceError("service.user.error.deleted_user_cannot_be_activated")
	ErrBannedUserCannotBeActivated  = newUserServiceError("service.user.error.banned_user_cannot_be_activated")
	ErrInvalidTimestamp             = newUserServiceError("service.user.error.invalid_timestamp")
)
