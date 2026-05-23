package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserDataDeletionTasksInboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return newCommonEventRepositoryWithClientID(conn, "user_data_deletion_tasks_inbox")
}
