package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"database/sql"
)

// insert chat to db
func insertChatToDB(tx *sql.Tx, chat *Chat) error {
	err := tx.QueryRow(`INSERT INTO chats (name) VALUES ($1) RETURNING id`, chat.Name).Scan(&chat.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "inserting chat to db", chat, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update chat in db
func updateChatInDB(tx *sql.Tx, chat *Chat) error {
	_, err := tx.Exec(`UPDATE chats SET name = $1, is_deleted = $2 WHERE id = $3`, chat.Name, chat.IsDeleted, chat.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "updating chat in db", chat, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// get chat from db by id
func getChatFromDB(id uint64) (*Chat, error) {
	db := pgDB.GetDB()

	var chat Chat

	err := db.QueryRow(`SELECT id, name, is_deleted FROM chats WHERE id = $1 AND is_deleted = FALSE`, id).Scan(&chat.ID, &chat.Name, &chat.IsDeleted)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("chat not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get chat from db by id", id, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return &chat, nil
}

// get private chat from db by members ids
func getPrivateChatFromDB(member1, member2 uint64) (*Chat, error) {
	var chat Chat
	db := pgDB.GetDB()
	err := db.QueryRow(`
		SELECT c.id, c.name, c.is_deleted FROM chats c
		JOIN chat_members cm ON cm.chat_id = c.id
		WHERE c.name IS NULL 
		AND ((added_by = $1 AND cm.user_id = $2) OR (added_by = $2 AND cm.user_id = $1))
	`, member1, member2).Scan(&chat.ID, &chat.Name, &chat.IsDeleted)
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.NotFound("private chat not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get private chat from db by members ids", map[string]interface{}{"member1": member1, "member2": member2}, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	chat.IsGroupChat = false

	return &chat, nil
}

// get chat list from db
func getChatListFromDB(userID uint64, isGroup bool, limit, offset int) ([]Chat, []*uint64, error) {
	db := pgDB.GetDB()

	query := `
		WITH user_chats AS (
			SELECT cm.chat_id, cm.added_at, COALESCE(cm.removed_at, NOW()) AS removed_at
			FROM chat_members cm
			WHERE cm.user_id = $1
		),
		last_messages AS (
			SELECT
				m.chat_id,
				MAX(m.sent_at) AS last_message_time,
				MAX(m.id) AS last_message_id
			FROM messages m
			JOIN user_chats uc ON m.chat_id = uc.chat_id
			WHERE 
				m.is_deleted = FALSE
				AND m.sent_at BETWEEN uc.added_at AND uc.removed_at
			GROUP BY m.chat_id
		)
		SELECT
			c.id, c.name, c.is_deleted, lm.last_message_id
		FROM chats c
		JOIN user_chats uc ON c.id = uc.chat_id
		LEFT JOIN last_messages lm ON c.id = lm.chat_id
		WHERE c.is_deleted = FALSE
	`

	if isGroup {
		query += ` AND c.name IS NOT NULL`
	} else {
		query += ` AND c.name IS NULL`
	}

	query += ` ORDER BY lm.last_message_time DESC NULLS LAST LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, userID, limit, offset)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil, appErr.NotFound("chats not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get chat list from sql rows", rows, err)
		return nil, nil, appErr.InternalServerError("internal server error")
	}

	return parseChatListFromSQLRows(rows)
}

// parse chats and message ids from sql rows
func parseChatListFromSQLRows(rows *sql.Rows) ([]Chat, []*uint64, error) {
	var chats []Chat
	var messageIDs []*uint64

	for rows.Next() {
		var chat Chat
		var messageID *uint64
		err := rows.Scan(&chat.ID, &chat.Name, &chat.IsDeleted, &messageID)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "creating chats from sql rows", rows, err)
			return nil, nil, appErr.InternalServerError("internal server error")
		}

		if chat.Name != nil {
			chat.IsGroupChat = true
		}

		chats = append(chats, chat)
		messageIDs = append(messageIDs, messageID)
	}

	if err := rows.Err(); err != nil {
		logger.GetInstance().Error(err.Error(), "creating chats from sql rows", rows, err)
		return nil, nil, appErr.InternalServerError("internal server error")
	}

	if len(chats) == 0 {
		return nil, nil, appErr.NotFound("chats not found")
	}

	return chats, messageIDs, nil
}
