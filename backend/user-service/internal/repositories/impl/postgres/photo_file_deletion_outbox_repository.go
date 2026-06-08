package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPhotoFileDeletionOutboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return newCommonEventRepositoryRepository(conn, "photo_file_deletion_outbox")
}
