package processing

type MediaCategory int16

const (
	CategoryImage MediaCategory = iota
	CategoryVideo
	CategoryAudio
)

func (c MediaCategory) String() string {
	switch c {
	case CategoryImage:
		return "image"
	case CategoryVideo:
		return "video"
	case CategoryAudio:
		return "audio"
	}
	return "unknown"
}

var allowedImageMimes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/gif":  {},
	"image/webp": {},
}

var allowedVideoMimes = map[string]struct{}{
	"video/mp4":        {},
	"video/webm":       {},
	"video/avi":        {},
	"video/quicktime":  {},
	"video/x-matroska": {},
}

var allowedAudioMimes = map[string]struct{}{
	"audio/mpeg": {},
	"audio/ogg":  {},
	"audio/wav":  {},
	"audio/webm": {},
	"audio/aac":  {},
	"audio/flac": {},
}

func (c MediaCategory) ValidateSourceMime(mimeType string) error {
	switch c {
	case CategoryImage:
		if _, ok := allowedImageMimes[mimeType]; !ok {
			return ErrMimeTypeNotSupported
		}
	case CategoryVideo:
		if _, ok := allowedVideoMimes[mimeType]; !ok {
			return ErrMimeTypeNotSupported
		}
	case CategoryAudio:
		if _, ok := allowedAudioMimes[mimeType]; !ok {
			return ErrMimeTypeNotSupported
		}
	default:
		return ErrMimeTypeNotSupported
	}
	return nil
}
