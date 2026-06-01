package filevariant

type FileVariantDomainError struct {
	message string
}

func newFileVariantDomainError(message string) *FileVariantDomainError {
	return &FileVariantDomainError{
		message: message,
	}
}

func (e *FileVariantDomainError) Error() string {
	return e.message
}

var (
	ErrUnknownVariantType = newFileVariantDomainError("domain.file_variant.unknown_variant_type")
)
