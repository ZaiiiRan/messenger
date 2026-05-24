package file

type FileError struct {
	message string
}

func NewFileError(message string) *FileError {
	return &FileError{message: message}
}

func (e *FileError) Error() string {
	return e.message
}
