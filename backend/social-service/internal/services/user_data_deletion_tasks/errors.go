package userdatadeletiontasksservice

type UserDataDeletionTasksServiceError struct {
	message string
}

func newUserDataDeletionTasksServiceError(message string) *UserDataDeletionTasksServiceError {
	return &UserDataDeletionTasksServiceError{
		message: message,
	}
}

func (e *UserDataDeletionTasksServiceError) Error() string {
	return e.message
}

var (
	ErrCreateUserDataDeletionTask   = newUserDataDeletionTasksServiceError("service.user_data_deletion_tasks.create_user_data_deletion_task_error")
	ErrProcessUserDataDeletionTasks = newUserDataDeletionTasksServiceError("service.user_data_deletion_tasks.process_user_data_deletion_tasks_error")
)
