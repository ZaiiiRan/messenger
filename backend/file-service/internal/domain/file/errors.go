package file

var (
	ErrUnknownFileStatus  = NewFileError("domain.file.error.unknown_file_status")
	ErrCannotUpdateStatus = NewFileError("domain.file.error.cannot_update_status")
	ErrUnknownFilePurpose = NewFileError("domain.file.error.unknown_file_purpose")
)
