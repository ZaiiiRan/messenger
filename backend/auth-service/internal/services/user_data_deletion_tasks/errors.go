package userdatadeletiontasksservice

type UserDataDeletionServiceError struct {
	message string
}

func newUserDataDeletionServiceError(message string) *UserDataDeletionServiceError {
	return &UserDataDeletionServiceError{
		message: message,
	}
}

func (e *UserDataDeletionServiceError) Error() string {
	return e.message
}

var (
	ErrCreateUserDataDeletionTask = newUserDataDeletionServiceError("service.user_deletion_tasks.create_user_data_deletion_task_error")
)
