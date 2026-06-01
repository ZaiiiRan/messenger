package mediaprocessor

import (
	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/processing"
)

type Request struct {
	ProcessingType processing.ProcessingType
	SourceMime     string
	SourceData     []byte
	SourceWidth    *int
	SourceHeight   *int
	SourceDuration *int
	Crop           *processing.CropParams
}

type ProcessedVariant struct {
	Variant  filevariant.VariantType
	Data     []byte
	MimeType string
	Width    *int
	Height   *int
	Duration *int
}

func intPtr(v int) *int { return &v }
