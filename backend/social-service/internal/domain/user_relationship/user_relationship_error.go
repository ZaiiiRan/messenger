package userrelationship

type UserRelationshipError struct {
	message string
}

func NewUserRelationshipError(message string) *UserRelationshipError {
	return &UserRelationshipError{
		message: message,
	}
}

func (e *UserRelationshipError) Error() string {
	return e.message
}
