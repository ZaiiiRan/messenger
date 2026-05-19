package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserDataDeletionTasksOutboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return &commonOutboxEventRepository{
		conn:      conn,
		tableName: "user_data_deletion_tasks_outbox",
	}
}
