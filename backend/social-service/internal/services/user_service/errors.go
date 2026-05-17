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

var ()
