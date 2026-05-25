package filedeletionfroms3tasksservice

type FileDeletionFromS3ServiceError struct {
	message string
}

func newFileDeletionFromS3ServiceError(message string) *FileDeletionFromS3ServiceError {
	return &FileDeletionFromS3ServiceError{
		message: message,
	}
}

func (e *FileDeletionFromS3ServiceError) Error() string {
	return e.message
}

var (
	ErrCreateFileDeletionFromS3Task   = newFileDeletionFromS3ServiceError("service.file_deletion_from_s3_tasks.create_file_deletion_from_s3_task_error")
	ErrProcessFileDeletionFromS3Tasks = newFileDeletionFromS3ServiceError("service.file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_error")
)
