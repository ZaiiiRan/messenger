package fileprocessingtasksservice

import (
	"context"
	"fmt"

	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/s3"
)

type fileStorage struct {
	userPhotosS3 *s3.S3Client
	chatFilesS3  *s3.S3Client
}

func newFileStorage(userPhotosS3, chatFilesS3 *s3.S3Client) *fileStorage {
	return &fileStorage{
		userPhotosS3: userPhotosS3,
		chatFilesS3:  chatFilesS3,
	}
}

func (s *fileStorage) clientFor(bucket filevariant.BucketType) (*s3.S3Client, error) {
	switch bucket {
	case filevariant.UserPhotosBucket:
		return s.userPhotosS3, nil
	case filevariant.ChatFilesBucket:
		return s.chatFilesS3, nil
	default:
		return nil, fmt.Errorf("unknown bucket type: %d", bucket)
	}
}

func (s *fileStorage) get(ctx context.Context, bucket filevariant.BucketType, key string) ([]byte, error) {
	client, err := s.clientFor(bucket)
	if err != nil {
		return nil, err
	}
	return client.Get(ctx, key)
}

func (s *fileStorage) put(ctx context.Context, bucket filevariant.BucketType, key string, data []byte, contentType string) error {
	client, err := s.clientFor(bucket)
	if err != nil {
		return err
	}
	return client.Put(ctx, key, data, contentType)
}
