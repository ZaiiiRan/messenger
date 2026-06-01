package mediaservice

import (
	"encoding/json"

	filepb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/file/v1"
	domainfile "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
)

func strPtr(s string) *string { return &s }

func extractOriginalVariant(f *filepb.File) (*filepb.FileVariant, bool) {
	for _, v := range f.GetVariants() {
		if v.GetVariantType() == filevariant.VariantOriginal.String() {
			return v, true
		}
	}
	return nil, false
}

func parseMeta(raw string) filevariant.Meta {
	var meta filevariant.Meta
	if raw == "" {
		return meta
	}
	_ = json.Unmarshal([]byte(raw), &meta)
	return meta
}

func buildProcessingStatusRequest(fileIDs []string) *filepb.UpdateFilesRequest {
	files := make([]*filepb.UpdateFileRequest, 0, len(fileIDs))
	for _, id := range fileIDs {
		files = append(files, &filepb.UpdateFileRequest{
			Id:     id,
			Status: strPtr(domainfile.FileStatusProcessing.String()),
		})
	}
	return &filepb.UpdateFilesRequest{Files: files}
}
