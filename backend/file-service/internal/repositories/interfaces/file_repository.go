package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
)

type FileRepository interface {
	CreateFiles(ctx context.Context, files []*file.File) error
	UpdateFiles(ctx context.Context, files []*file.File) error
	DeleteFiles(ctx context.Context, files []*file.File) error
	QueryFiles(ctx context.Context, query *models.QueryFilesDal, locked bool, skipLocked bool) ([]*file.File, error)
}
