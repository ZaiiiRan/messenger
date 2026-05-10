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
	ErrCreateUserDataDeletionTasks  = newUserDataDeletionServiceError("service.user_data_deletion_tasks.create_user_data_deletion_tasks_error")
	ErrProcessUserDataDeletionTasks = newUserDataDeletionServiceError("service.user_data_deletion_tasks.process_user_data_deletion_tasks_error")
)
