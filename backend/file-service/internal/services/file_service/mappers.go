package fileservice

import (
	"encoding/json"

	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/utils"
)

func fileVariantToProto(fv *filevariant.FileVariant, urlForDownload *string, bucketStr *string, includeProprietaryInfo bool) (*pb.FileVariant, error) {
	meta, err := json.Marshal(fv.GetMeta())
	if err != nil {
		return nil, err
	}

	var s3Key, bucket, createdAt, updatedAt *string
	if includeProprietaryInfo {
		s3Key = utils.StringPtr(fv.GetS3Key())
		bucket = bucketStr
		createdAt = utils.StringPtr(utils.FormatTimestamp(fv.GetCreatedAt()))
		updatedAt = utils.StringPtr(utils.FormatTimestamp(fv.GetUpdatedAt()))
	}

	return &pb.FileVariant{
		Id:             fv.GetID(),
		VariantType:    fv.GetVariant().String(),
		Size:           fv.GetSize(),
		MimeType:       fv.GetMimeType(),
		Meta:           string(meta),
		UrlForDownload: urlForDownload,
		S3Key:          s3Key,
		Bucket:         bucket,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

func fileToProto(f *file.File, variants []*pb.FileVariant, includeProprietaryInfo bool) *pb.File {
	var createdAt, updatedAt *string
	if includeProprietaryInfo {
		createdAt = utils.StringPtr(utils.FormatTimestamp(f.GetCreatedAt()))
		updatedAt = utils.StringPtr(utils.FormatTimestamp(f.GetUpdatedAt()))
	}

	return &pb.File{
		Id:        f.GetID(),
		OwnerId:   f.GetOwnerID(),
		Status:    f.GetStatus().String(),
		Variants:  variants,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
