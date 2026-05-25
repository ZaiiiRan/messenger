package filedeletiontasksservice

type FileDeletionTasksServiceError struct {
	message string
}

func newFileDeletionTasksServiceError(message string) *FileDeletionTasksServiceError {
	return &FileDeletionTasksServiceError{
		message: message,
	}
}

func (e *FileDeletionTasksServiceError) Error() string {
	return e.message
}

var (
	ErrCreateFileDeletionTask   = newFileDeletionTasksServiceError("service.file_deletion_tasks.create_file_deletion_task_error")
	ErrProcessFileDeletionTasks = newFileDeletionTasksServiceError("service.file_deletion_tasks.process_file_deletion_tasks_error")
)
