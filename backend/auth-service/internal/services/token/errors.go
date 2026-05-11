package tokenservice

type TokenServiceError struct {
	message string
}

func newTokenServiceError(message string) *TokenServiceError {
	return &TokenServiceError{message: message}
}

func (e *TokenServiceError) Error() string {
	return e.message
}

var (
	ErrUserVersionOrExistedRefreshTokenNotProvided = newTokenServiceError("user version or existed refresh token is not provided")
)
