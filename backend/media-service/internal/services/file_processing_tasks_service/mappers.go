package fileprocessingtasksservice

import (
	"encoding/json"
	"fmt"

	filepb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/file/v1"
	domainfile "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
)

type uploadedVariant struct {
	variant  filevariant.VariantType
	s3Key    string
	mimeType string
	meta     filevariant.Meta
}

func strPtr(s string) *string { return &s }

func buildFinalizeRequest(task *TaskPayload, variants []uploadedVariant) (*filepb.UpdateFilesRequest, error) {
	fileReq := &filepb.UpdateFileRequest{
		Id:          task.FileID,
		FilePurpose: strPtr(task.Purpose),
		Status:      strPtr(domainfile.FileStatusReady.String()),
	}

	keepsOriginal := false
	for _, v := range variants {
		if v.variant == filevariant.VariantOriginal {
			keepsOriginal = true
			continue
		}
		metaBytes, err := json.Marshal(v.meta)
		if err != nil {
			return nil, fmt.Errorf("marshal variant meta: %w", err)
		}
		fileReq.AddVariants = append(fileReq.AddVariants, &filepb.AddVariantRequest{
			VariantType: v.variant.String(),
			S3Key:       v.s3Key,
			Meta:        string(metaBytes),
		})
	}

	if !keepsOriginal && task.OriginalVariantID != "" {
		fileReq.DeleteVariants = []string{task.OriginalVariantID}
	}

	if len(fileReq.AddVariants) == 0 {
		fileReq.FilePurpose = nil
	}

	return &filepb.UpdateFilesRequest{Files: []*filepb.UpdateFileRequest{fileReq}}, nil
}

func variantKey(fileID string, variant filevariant.VariantType, mimeType string) string {
	return fmt.Sprintf("%s/%s%s", fileID, variant.String(), extForMime(mimeType))
}

func extForMime(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "audio/mpeg":
		return ".mp3"
	case "audio/aac":
		return ".aac"
	case "audio/ogg":
		return ".ogg"
	case "audio/webm":
		return ".weba"
	case "audio/flac":
		return ".flac"
	case "audio/wav":
		return ".wav"
	default:
		return ""
	}
}
