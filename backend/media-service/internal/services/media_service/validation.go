package mediaservice

import (
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	filepb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/file/v1"
	mediapb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/media/v1"
	domainfile "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/processing"
	fileprocessingtasksservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_tasks_service"
)

func validateAndBuildPayloads(
	req *mediapb.ProcessFilesRequest,
	filesByID map[string]*filepb.File,
	claims *jwt.UserClaims,
) ([]*fileprocessingtasksservice.TaskPayload, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)
	payloads := make([]*fileprocessingtasksservice.TaskPayload, 0, len(req.Files))
	seen := make(map[string]struct{}, len(req.Files))

	for i, fr := range req.Files {
		fileID := fr.GetFileId()
		if fileID == "" {
			verr[fmt.Sprintf("files[%d].file_id", i)] = ErrFileIDRequired.Error()
			continue
		}
		if _, dup := seen[fileID]; dup {
			verr[fmt.Sprintf("files[%d].file_id", i)] = ErrDuplicateFile.Error()
			continue
		}
		seen[fileID] = struct{}{}

		pt, err := processing.ProcessingTypeFromString(fr.GetProcessingType())
		if err != nil {
			verr[fmt.Sprintf("files[%d].processing_type", i)] = err.Error()
			continue
		}

		var crop *processing.CropParams
		if pt.RequiresCrop() {
			ap := fr.GetAvatarParams()
			if ap == nil {
				verr[fmt.Sprintf("files[%d].avatar_params", i)] = processing.ErrCropParamsRequired.Error()
				continue
			}
			crop = &processing.CropParams{
				X:      int(ap.GetCropX()),
				Y:      int(ap.GetCropY()),
				Width:  int(ap.GetCropWidth()),
				Height: int(ap.GetCropHeight()),
			}
		} else if fr.GetAvatarParams() != nil {
			verr[fmt.Sprintf("files[%d].avatar_params", i)] = processing.ErrCropParamsNotAllowed.Error()
			continue
		}

		f := filesByID[fileID]
		if f == nil {
			verr[fmt.Sprintf("files[%d].file_id", i)] = ErrFileNotFound.Error()
			continue
		}
		if f.GetOwnerId() != claims.Id {
			verr[fmt.Sprintf("files[%d].file_id", i)] = ErrFileNotOwned.Error()
			continue
		}
		if f.GetStatus() != domainfile.FileStatusUploaded.String() {
			verr[fmt.Sprintf("files[%d].file_id", i)] = ErrFileNotUploaded.Error()
			continue
		}

		original, ok := extractOriginalVariant(f)
		if !ok {
			verr[fmt.Sprintf("files[%d].file_id", i)] = ErrOriginalNotFound.Error()
			continue
		}
		mime := original.GetMimeType()
		meta := parseMeta(original.GetMeta())

		category := pt.Category()
		if err := category.ValidateSourceMime(mime); err != nil {
			verr[fmt.Sprintf("files[%d].processing_type", i)] = err.Error()
			continue
		}

		switch category {
		case processing.CategoryImage:
			if err := processing.ValidateImageDimensions(meta.Width, meta.Height); err != nil {
				verr[fmt.Sprintf("files[%d].file_id", i)] = err.Error()
				continue
			}
			if pt.RequiresCrop() {
				if err := processing.ValidateCrop(crop, meta.Width, meta.Height); err != nil {
					verr[fmt.Sprintf("files[%d].avatar_params", i)] = err.Error()
					continue
				}
			}
		case processing.CategoryVideo:
			if err := processing.ValidateVideoDimensions(meta.Width, meta.Height); err != nil {
				verr[fmt.Sprintf("files[%d].file_id", i)] = err.Error()
				continue
			}
			if err := processing.ValidateDuration(meta.Duration, processing.MaxVideoDurationSeconds); err != nil {
				verr[fmt.Sprintf("files[%d].file_id", i)] = err.Error()
				continue
			}
		case processing.CategoryAudio:
			if err := processing.ValidateDuration(meta.Duration, processing.MaxAudioDurationSeconds); err != nil {
				verr[fmt.Sprintf("files[%d].file_id", i)] = err.Error()
				continue
			}
		}

		payloads = append(payloads, &fileprocessingtasksservice.TaskPayload{
			FileID:            fileID,
			OwnerID:           f.GetOwnerId(),
			ProcessingType:    pt.String(),
			Purpose:           pt.Purpose().String(),
			Bucket:            int16(pt.Bucket()),
			OriginalVariantID: original.GetId(),
			OriginalS3Key:     original.GetS3Key(),
			OriginalMime:      mime,
			Width:             meta.Width,
			Height:            meta.Height,
			Duration:          meta.Duration,
			Crop:              crop,
		})
	}

	return payloads, verr
}
