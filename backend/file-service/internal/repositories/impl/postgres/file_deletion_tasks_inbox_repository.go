package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewFileDeletionTasksInboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return newCommonEventRepositoryWithClientID(conn, "file_deletion_tasks_inbox")
}
