package fileservice

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

func (s *fileStorage) userPhotoPresignedUpload(
	ctx context.Context,
	key string,
	contentType string,
	maxSize int64,
) (string, map[string]string, error) {
	return s.userPhotosS3.PresignUpload(ctx, key, contentType, maxSize)
}

func (s *fileStorage) userPhotoPresignedDownload(ctx context.Context, key string) (string, error) {
	return s.userPhotosS3.PresignDownload(ctx, key)
}

func (s *fileStorage) userPhotoStat(ctx context.Context, key string) (*s3.ObjectInfo, error) {
	return s.userPhotosS3.Stat(ctx, key)
}

func (s *fileStorage) chatFilePresignedUpload(
	ctx context.Context,
	key string,
	contentType string,
	maxSize int64,
) (string, map[string]string, error) {
	return s.chatFilesS3.PresignUpload(ctx, key, contentType, maxSize)
}

func (s *fileStorage) chatFilePresignedDownload(ctx context.Context, key string) (string, error) {
	return s.chatFilesS3.PresignDownload(ctx, key)
}

func (s *fileStorage) chatFileStat(ctx context.Context, key string) (*s3.ObjectInfo, error) {
	return s.chatFilesS3.Stat(ctx, key)
}
