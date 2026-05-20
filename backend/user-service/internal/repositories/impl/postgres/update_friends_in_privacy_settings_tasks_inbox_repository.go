package postgresimpl

import (
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUpdateFriendsInPrivacySettingsTasksInboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return newCommonEventRepositoryRepository(conn, "update_friends_in_privacy_settings_inbox")
}
