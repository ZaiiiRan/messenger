package filevariant

type FileVariantError struct {
	message string
}

func NewFileVariantError(message string) *FileVariantError {
	return &FileVariantError{message: message}
}

func (e *FileVariantError) Error() string {
	return e.message
}
