package fileservice

type FileServiceError struct {
	message string
}

func newFileServiceError(message string) *FileServiceError {
	return &FileServiceError{
		message: message,
	}
}

func (e *FileServiceError) Error() string {
	return e.message
}

var (
	ErrFileIdsIsRequired = newFileServiceError("service.file.error.file_ids_is_required")
	ErrTooManyFiles      = newFileServiceError("service.file.error.too_many_files")
)
