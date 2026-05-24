package filevariant

import (
	"regexp"

	file "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
)

var mimeTypeRe = regexp.MustCompile(`^[a-z]+/[a-z0-9][a-z0-9!#$&\-^_.+]*$`)

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

func ValidateMeta(meta Meta, purpose file.FilePurpose) error {
	switch purpose {
	case file.PurposeUserPhoto, file.PurposeChatPhoto, file.PurposeChatAvatar:
		if meta.Width == nil || meta.Height == nil || *meta.Width <= 0 || *meta.Height <= 0 {
			return ErrInvalidMeta
		}
	case file.PurposeChatVideo:
		if meta.Width == nil || meta.Height == nil || *meta.Width <= 0 || *meta.Height <= 0 {
			return ErrInvalidMeta
		}
		if meta.Duration == nil || *meta.Duration <= 0 {
			return ErrInvalidMeta
		}
	case file.PurposeChatAudio:
		if meta.Duration == nil || *meta.Duration <= 0 {
			return ErrInvalidMeta
		}
	}
	return nil
}

var allowedVariants = map[file.FilePurpose]map[VariantType]struct{}{
	file.PurposeUserPhoto: {
		VariantOriginal:     {},
		VariantAvatar:       {},
		VariantAvatarSmall:  {},
		VariantAvatarMiddle: {},
	},
	file.PurposeChatAvatar: {
		VariantOriginal:     {},
		VariantAvatar:       {},
		VariantAvatarSmall:  {},
		VariantAvatarMiddle: {},
	},
	file.PurposeChatPhoto: {
		VariantOriginal:    {},
		VariantCompressed:  {},
		VariantThumb:       {},
	},
	file.PurposeChatVideo: {
		VariantOriginal:    {},
		VariantCompressed:  {},
		Variant360p:        {},
		Variant720p:        {},
		Variant1080p:       {},
		VariantThumb:       {},
	},
	file.PurposeChatAudio: {
		VariantOriginal: {},
	},
	file.PurposeChatFile: {
		VariantOriginal: {},
	},
}

func ValidateVariantType(variant VariantType, purpose file.FilePurpose) error {
	allowed, ok := allowedVariants[purpose]
	if !ok {
		return ErrVariantNotAllowedForPurpose
	}
	if _, ok := allowed[variant]; !ok {
		return ErrVariantNotAllowedForPurpose
	}
	return nil
}

func ValidateMimeType(mimeType string, purpose file.FilePurpose) error {
	switch purpose {
	case file.PurposeUserPhoto, file.PurposeChatPhoto, file.PurposeChatAvatar:
		if _, ok := allowedImageMimes[mimeType]; !ok {
			return ErrInvalidMimeType
		}
	case file.PurposeChatVideo:
		if _, ok := allowedVideoMimes[mimeType]; !ok {
			return ErrInvalidMimeType
		}
	case file.PurposeChatAudio:
		if _, ok := allowedAudioMimes[mimeType]; !ok {
			return ErrInvalidMimeType
		}
	case file.PurposeChatFile:
		if !mimeTypeRe.MatchString(mimeType) {
			return ErrInvalidMimeType
		}
	}
	return nil
}
