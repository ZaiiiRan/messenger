package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewFileProcessingCompletedOutboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return newCommonEventRepository(conn, "file_processing_completed_outbox")
}
