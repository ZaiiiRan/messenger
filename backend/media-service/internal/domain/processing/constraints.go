package processing

const (
	MinImageDimension = 16
	MaxImageDimension = 12000

	MinVideoDimension = 16
	MaxVideoDimension = 8000

	MinCropDimension = 16

	MaxAudioDurationSeconds = 3600
	MaxVideoDurationSeconds = 3600
)

type CropParams struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func ValidateImageDimensions(width, height *int) error {
	if width == nil || height == nil {
		return ErrMissingDimensions
	}
	if *width < MinImageDimension || *height < MinImageDimension ||
		*width > MaxImageDimension || *height > MaxImageDimension {
		return ErrInvalidImageDimensions
	}
	return nil
}

func ValidateVideoDimensions(width, height *int) error {
	if width == nil || height == nil {
		return ErrMissingDimensions
	}
	if *width < MinVideoDimension || *height < MinVideoDimension ||
		*width > MaxVideoDimension || *height > MaxVideoDimension {
		return ErrInvalidVideoDimensions
	}
	return nil
}

func ValidateDuration(duration *int, maxSeconds int) error {
	if duration == nil {
		return ErrMissingDuration
	}
	if *duration <= 0 {
		return ErrInvalidDuration
	}
	if *duration > maxSeconds {
		return ErrDurationTooLong
	}
	return nil
}

func ValidateCrop(crop *CropParams, srcWidth, srcHeight *int) error {
	if crop == nil {
		return ErrCropParamsRequired
	}
	if srcWidth == nil || srcHeight == nil {
		return ErrMissingDimensions
	}
	if crop.Width < MinCropDimension || crop.Height < MinCropDimension {
		return ErrInvalidCropParams
	}
	if crop.X < 0 || crop.Y < 0 {
		return ErrInvalidCropParams
	}
	if crop.X+crop.Width > *srcWidth || crop.Y+crop.Height > *srcHeight {
		return ErrCropOutOfBounds
	}
	return nil
}
