package models

import "time"

type FileProcessingCompletedVariant struct {
	VariantType string `json:"variant_type"`
	MimeType    string `json:"mime_type"`
	Width       *int   `json:"width,omitempty"`
	Height      *int   `json:"height,omitempty"`
	Duration    *int   `json:"duration,omitempty"`
}

type FileProcessingCompletedEvent struct {
	Id             string                           `json:"id"`
	FileId         string                           `json:"file_id"`
	OwnerId        string                           `json:"owner_id"`
	ProcessingType string                           `json:"processing_type"`
	Purpose        string                           `json:"purpose"`
	Status         string                           `json:"status"`
	Variants       []FileProcessingCompletedVariant `json:"variants"`
	CompletedAt    time.Time                        `json:"completed_at"`
}
