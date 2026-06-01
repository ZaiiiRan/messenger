package file

type FileDomainError struct {
	message string
}

func newFileDomainError(message string) *FileDomainError {
	return &FileDomainError{
		message: message,
	}
}

func (e *FileDomainError) Error() string {
	return e.message
}

var (
	ErrUnknownFileStatus  = newFileDomainError("domain.file.unknown_file_status")
	ErrUnknownFilePurpose = newFileDomainError("domain.file.unknown_file_purpose")
)
