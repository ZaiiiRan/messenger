package file

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/utils"
)

type File struct {
	id        string
	ownerID   string
	status    int16
	createdAt time.Time
	updatedAt time.Time
}

func New(
	ownerID string,
) *File {
	now := time.Now()
	return &File{
		ownerID:   ownerID,
		status:    int16(FileStatusWaitingForUpload),
		createdAt: now,
		updatedAt: now,
	}
}

func FromStorage(
	id string,
	ownerID string,
	status int16,
	createdAt time.Time,
	updatedAt time.Time,
) *File {
	return &File{
		id:        id,
		ownerID:   ownerID,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (f *File) GetID() string           { return f.id }
func (f *File) GetOwnerID() string      { return f.ownerID }
func (f *File) GetStatus() FileStatus   { return FileStatus(f.status) }
func (f *File) GetCreatedAt() time.Time { return f.createdAt }
func (f *File) GetUpdatedAt() time.Time { return f.updatedAt }

func (f *File) SetID(id string) {
	if f.id == "" {
		f.id = id
	}
}

func (f *File) SetStatus(status FileStatus) error {
	convertedStatus := int16(status)

	switch status {
	case FileStatusWaitingForUpload:
		return ErrCannotUpdateStatus
	case FileStatusUploaded:
		if FileStatus(f.status) != FileStatusWaitingForUpload {
			return ErrCannotUpdateStatus
		}
	case FileStatusProcessing:
		if FileStatus(f.status) != FileStatusUploaded {
			return ErrCannotUpdateStatus
		}
	case FileStatusReady:
		if FileStatus(f.status) != FileStatusProcessing && FileStatus(f.status) != FileStatusUploaded {
			return ErrCannotUpdateStatus
		}
	case FileStatusInUse:
		if FileStatus(f.status) != FileStatusReady && FileStatus(f.status) != FileStatusUploaded {
			return ErrCannotUpdateStatus
		}
	default:
		return ErrUnknownFileStatus
	}
	f.status = convertedStatus
	return nil
}

func (f *File) SetUpdatedAt(now *time.Time) {
	if now == nil {
		now = utils.TimePtr(time.Now())
	}
	f.updatedAt = *now
}
