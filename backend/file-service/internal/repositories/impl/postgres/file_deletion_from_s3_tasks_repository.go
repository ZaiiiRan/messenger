package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewFileDeletionFromS3TasksRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return newCommonEventRepository(conn, "file_deletion_from_s3_tasks")
}
