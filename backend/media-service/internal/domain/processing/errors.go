package processing

type ProcessingDomainError struct {
	message string
}

func newProcessingDomainError(message string) *ProcessingDomainError {
	return &ProcessingDomainError{
		message: message,
	}
}

func (e *ProcessingDomainError) Error() string {
	return e.message
}

var (
	ErrUnknownProcessingType  = newProcessingDomainError("domain.processing.unknown_processing_type")
	ErrCropParamsRequired     = newProcessingDomainError("domain.processing.crop_params_required")
	ErrCropParamsNotAllowed   = newProcessingDomainError("domain.processing.crop_params_not_allowed")
	ErrInvalidCropParams      = newProcessingDomainError("domain.processing.invalid_crop_params")
	ErrCropOutOfBounds        = newProcessingDomainError("domain.processing.crop_out_of_bounds")
	ErrInvalidImageDimensions = newProcessingDomainError("domain.processing.invalid_image_dimensions")
	ErrInvalidVideoDimensions = newProcessingDomainError("domain.processing.invalid_video_dimensions")
	ErrInvalidDuration        = newProcessingDomainError("domain.processing.invalid_duration")
	ErrDurationTooLong        = newProcessingDomainError("domain.processing.duration_too_long")
	ErrMimeTypeNotSupported   = newProcessingDomainError("domain.processing.mime_type_not_supported")
	ErrMissingDimensions      = newProcessingDomainError("domain.processing.missing_dimensions")
	ErrMissingDuration        = newProcessingDomainError("domain.processing.missing_duration")
)
