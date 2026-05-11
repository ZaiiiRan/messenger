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
	ErrMarshalPayload             = newUserDataDeletionServiceError("service.user_deletion_tasks.payload_marshal_error")
	ErrCreateUserDataDeletionTask = newUserDataDeletionServiceError("service.user_deletion_tasks.create_user_data_deletion_task_error")
	ErrSendUserDataDeletionTasks  = newUserDataDeletionServiceError("service.user_deletion_tasks.send_user_data_deletion_tasks_error")
)
