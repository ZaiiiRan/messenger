package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserRelationshipChangesTasksOutboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return &commonEventRepository{
		tableName: "user_relationship_changes_tasks_outbox",
		conn:      conn,
	}
}
