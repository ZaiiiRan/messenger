package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
)

type FileCacheRepository interface {
	GetFileByID(ctx context.Context, id string) (*file.File, error)
	SetFileByID(ctx context.Context, file *file.File) error
	DelFileByID(ctx context.Context, id string) error
	GetFileList(ctx context.Context, query *models.QueryFilesDal) ([]*file.File, error)
	SetFileList(ctx context.Context, query *models.QueryFilesDal, files []*file.File) error
	InvalidateFileList(ctx context.Context, query *models.QueryFilesDal) error
	InvalidateFileListsByOwnerID(ctx context.Context, ownerID string) error
}
