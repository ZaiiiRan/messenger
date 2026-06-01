package fileprocessingtasksservice

import "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/processing"

type TaskPayload struct {
	FileID            string                 `json:"file_id"`
	OwnerID           string                 `json:"owner_id"`
	ProcessingType    string                 `json:"processing_type"`
	Purpose           string                 `json:"purpose"`
	Bucket            int16                  `json:"bucket"`
	OriginalVariantID string                 `json:"original_variant_id"`
	OriginalS3Key     string                 `json:"original_s3_key"`
	OriginalMime      string                 `json:"original_mime"`
	Width             *int                   `json:"width,omitempty"`
	Height            *int                   `json:"height,omitempty"`
	Duration          *int                   `json:"duration,omitempty"`
	Crop              *processing.CropParams `json:"crop,omitempty"`
}
