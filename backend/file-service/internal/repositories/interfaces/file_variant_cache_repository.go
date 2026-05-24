package interfaces

import (
	"context"

	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
)

type FileVariantCacheRepository interface {
	GetFileVariantsByFileId(ctx context.Context, fileId string) ([]*filevariant.FileVariant, error)
	SetFileVariantsByFileId(ctx context.Context, fileId string, variants []*filevariant.FileVariant) error
	DelFileVariantsByFileId(ctx context.Context, fileId string) error
	GetFileVariantsList(ctx context.Context, query *models.QueryFileVariantsDal) ([]*filevariant.FileVariant, error)
	SetFileVariantsList(ctx context.Context, query *models.QueryFileVariantsDal, fileVariants []*filevariant.FileVariant) error
	InvalidateFileVariantsList(ctx context.Context, query *models.QueryFileVariantsDal) error
	InvalidateFileVariantListsByFileID(ctx context.Context, fileID string) error
}
