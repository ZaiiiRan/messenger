package fileprocessingtasksservice

type FileProcessingTasksServiceError struct {
	message string
}

func newFileProcessingTasksServiceError(message string) *FileProcessingTasksServiceError {
	return &FileProcessingTasksServiceError{
		message: message,
	}
}

func (e *FileProcessingTasksServiceError) Error() string {
	return e.message
}

var (
	ErrCreateFileProcessingTasks  = newFileProcessingTasksServiceError("service.file_processing_tasks.create_file_processing_tasks_error")
	ErrProcessFileProcessingTasks = newFileProcessingTasksServiceError("service.file_processing_tasks.process_file_processing_tasks_error")
)
