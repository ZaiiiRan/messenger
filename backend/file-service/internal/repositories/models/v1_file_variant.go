package models

import (
	"encoding/json"
	"time"

	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
)

type V1FileVariantDal struct {
	Id        string          `json:"id" db:"id"`
	FileId    string          `json:"file_id" db:"file_id"`
	Variant   int16           `json:"variant" db:"variant"`
	Bucket    int16           `json:"bucket" db:"bucket"`
	S3Key     string          `json:"s3_key" db:"s3_key"`
	MimeType  string          `json:"mime_type" db:"mime_type"`
	Size      int64           `json:"size" db:"size"`
	Meta      json.RawMessage `json:"meta" db:"meta"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

func V1FileVariantFromDomain(fv *filevariant.FileVariant) (V1FileVariantDal, error) {
	if fv == nil {
		return V1FileVariantDal{}, nil
	}

	metaBytes, err := json.Marshal(fv.GetMeta())
	if err != nil {
		return V1FileVariantDal{}, err
	}

	return V1FileVariantDal{
		Id:        fv.GetID(),
		FileId:    fv.GetFileID(),
		Variant:   int16(fv.GetVariant()),
		Bucket:    int16(fv.GetBucket()),
		S3Key:     fv.GetS3Key(),
		MimeType:  fv.GetMimeType(),
		Size:      fv.GetSize(),
		Meta:      metaBytes,
		CreatedAt: fv.GetCreatedAt(),
		UpdatedAt: fv.GetUpdatedAt(),
	}, nil
}

func (v V1FileVariantDal) IsNull() bool { return false }
func (v V1FileVariantDal) Index(i int) any {
	switch i {
	case 0:
		return v.Id
	case 1:
		return v.FileId
	case 2:
		return v.Variant
	case 3:
		return v.Bucket
	case 4:
		return v.S3Key
	case 5:
		return v.MimeType
	case 6:
		return v.Size
	case 7:
		return v.Meta
	case 8:
		return v.CreatedAt
	case 9:
		return v.UpdatedAt
	default:
		return nil
	}
}

func (v V1FileVariantDal) ToDomain() (*filevariant.FileVariant, error) {
	return filevariant.FromStorage(
		v.Id, v.FileId,
		v.Variant, v.Bucket,
		v.S3Key,
		v.MimeType,
		v.Size,
		v.Meta,
		v.CreatedAt,
		v.UpdatedAt,
	)
}
