package interfaces

import (
	"context"

	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
)

type FileVariantRepository interface {
	CreateFileVariants(ctx context.Context, fileVariants []*filevariant.FileVariant) error
	UpdateFileVariants(ctx context.Context, fileVariants []*filevariant.FileVariant) error
	DeleteFileVariants(ctx context.Context, fileVariants []*filevariant.FileVariant) error
	QueryFileVariants(ctx context.Context, query *models.QueryFileVariantsDal, locked bool) ([]*filevariant.FileVariant, error)
}
