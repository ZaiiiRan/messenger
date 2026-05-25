package filevariant

import (
	"encoding/json"
	"time"

	file "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/utils"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
)

type FileVariant struct {
	id        string
	fileID    string
	variant   int16
	bucket    int16
	s3Key     string
	mimeType  string
	size      int64
	meta      Meta
	createdAt time.Time
	updatedAt time.Time
}

func New(fileId string, bucket BucketType, s3Key string, mimeType string) (*FileVariant, error) {
	now := time.Now()

	return &FileVariant{
		fileID:    fileId,
		variant:   int16(VariantOriginal),
		bucket:    int16(bucket),
		s3Key:     s3Key,
		mimeType:  mimeType,
		size:      0,
		meta:      Meta{},
		createdAt: now,
		updatedAt: now,
	}, nil
}

func NewCustom(
	fileId string,
	variant VariantType,
	bucket BucketType,
	s3Key string,
	mimeType string,
	size int64,
	meta json.RawMessage,
	purpose file.FilePurpose,
) (*FileVariant, validationerror.ValidationError) {
	now := time.Now()
	fv := &FileVariant{
		fileID:    fileId,
		variant:   int16(variant),
		bucket:    int16(bucket),
		s3Key:     s3Key,
		createdAt: now,
		updatedAt: now,
	}

	verr := make(validationerror.ValidationError)

	if err := ValidateVariantType(variant, purpose); err != nil {
		verr["variant"] = err.Error()
	}

	if err := fv.SetMimeType(mimeType, purpose); err != nil {
		verr["mime_type"] = err.Error()
	}

	if err := fv.SetSize(size); err != nil {
		verr["size"] = err.Error()
	}

	var metaStruct Meta
	if err := json.Unmarshal(meta, &metaStruct); err != nil {
		verr["meta"] = ErrUmarshalMeta.Error()
	} else if err := fv.SetMeta(metaStruct, purpose); err != nil {
		verr["meta"] = err.Error()
	}

	if len(verr) > 0 {
		return nil, verr
	}
	return fv, nil
}

func FromStorage(
	id string,
	fileId string,
	variant int16,
	bucket int16,
	s3Key string,
	mimeType string,
	size int64,
	metaBytes json.RawMessage,
	createdAt time.Time,
	updatedAt time.Time,
) (*FileVariant, error) {
	meta := Meta{}
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return nil, ErrUmarshalMeta
	}

	return &FileVariant{
		id:        id,
		fileID:    fileId,
		variant:   variant,
		bucket:    bucket,
		s3Key:     s3Key,
		mimeType:  mimeType,
		size:      size,
		meta:      meta,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (fv *FileVariant) GetID() string           { return fv.id }
func (fv *FileVariant) GetFileID() string       { return fv.fileID }
func (fv *FileVariant) GetVariant() VariantType { return VariantType(fv.variant) }
func (fv *FileVariant) GetBucket() BucketType   { return BucketType(fv.bucket) }
func (fv *FileVariant) GetS3Key() string        { return fv.s3Key }
func (fv *FileVariant) GetMimeType() string     { return fv.mimeType }
func (fv *FileVariant) GetSize() int64          { return fv.size }
func (fv *FileVariant) GetMeta() Meta           { return fv.meta }
func (fv *FileVariant) GetCreatedAt() time.Time { return fv.createdAt }
func (fv *FileVariant) GetUpdatedAt() time.Time { return fv.updatedAt }

func (fv *FileVariant) SetID(id string) {
	if fv.id == "" {
		fv.id = id
	}
}

func (fv *FileVariant) SetFileID(fileID string) {
	if fv.fileID == "" {
		fv.fileID = fileID
	}
}

func (fv *FileVariant) SetMimeType(mimeType string, purpose file.FilePurpose) error {
	if err := ValidateMimeType(mimeType, purpose); err != nil {
		return err
	}
	fv.mimeType = mimeType
	return nil
}

func (fv *FileVariant) SetSize(size int64) error {
	if size < 0 {
		return ErrInvalidSize
	}
	fv.size = size
	return nil
}

func (fv *FileVariant) SetVariant(variant VariantType, purpose file.FilePurpose) error {
	if err := ValidateVariantType(variant, purpose); err != nil {
		return err
	}
	fv.variant = int16(variant)
	return nil
}

func (fv *FileVariant) SetMeta(meta Meta, purpose file.FilePurpose) error {
	if err := ValidateMeta(meta, purpose); err != nil {
		return err
	}

	switch purpose {
	case file.PurposeUserPhoto, file.PurposeChatPhoto, file.PurposeChatAvatar:
		meta.Duration = nil
	case file.PurposeChatAudio:
		meta.Width = nil
		meta.Height = nil
	}

	fv.meta = meta
	return nil
}

func (fv *FileVariant) SetUpdatedAt(now *time.Time) {
	if now == nil {
		now = utils.TimePtr(time.Now())
	}
	fv.updatedAt = *now
}
