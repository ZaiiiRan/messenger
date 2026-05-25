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
	ErrFileIdsIsRequired     = newFileServiceError("service.file.error.file_ids_is_required")
	ErrVariantIdsIsRequired  = newFileServiceError("service.file.error.variant_ids_is_required")
	ErrTooManyFiles          = newFileServiceError("service.file.error.too_many_files")
	ErrFileNotFound          = newFileServiceError("service.file.error.file_not_found")
	ErrFilePurposeIsRequired = newFileServiceError("service.file.error.file_purpose_is_required")
	ErrS3KeyIsRequired       = newFileServiceError("service.file.error.s3_key_is_required")
	ErrS3ObjectNotFound      = newFileServiceError("service.file.error.s3_object_not_found")
)
