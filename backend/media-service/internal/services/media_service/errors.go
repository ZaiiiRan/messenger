package mediaservice

type MediaServiceError struct {
	message string
}

func newMediaServiceError(message string) *MediaServiceError {
	return &MediaServiceError{
		message: message,
	}
}

func (e *MediaServiceError) Error() string {
	return e.message
}

var (
	ErrProcessFiles = newMediaServiceError("service.media.process_files_error")

	ErrNoFiles          = newMediaServiceError("service.media.no_files")
	ErrTooManyFiles     = newMediaServiceError("service.media.too_many_files")
	ErrFileIDRequired   = newMediaServiceError("service.media.file_id_required")
	ErrFileNotFound     = newMediaServiceError("service.media.file_not_found")
	ErrFileNotOwned     = newMediaServiceError("service.media.file_not_owned")
	ErrFileNotUploaded  = newMediaServiceError("service.media.file_not_uploaded")
	ErrOriginalNotFound = newMediaServiceError("service.media.original_variant_not_found")
	ErrDuplicateFile    = newMediaServiceError("service.media.duplicate_file")
)

const maxFilesPerRequest = 100
