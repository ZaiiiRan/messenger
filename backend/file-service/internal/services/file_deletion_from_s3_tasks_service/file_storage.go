package filedeletionfroms3tasksservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/s3"
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

func (s *fileStorage) userPhotosDeleteMany(ctx context.Context, keys []string) error {
	return s.userPhotosS3.DeleteMany(ctx, keys)
}

func (s *fileStorage) chatFilesDeleteMany(ctx context.Context, keys []string) error {
	return s.chatFilesS3.DeleteMany(ctx, keys)
}
