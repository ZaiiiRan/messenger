package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"database/sql"
)

// insert chat to db
func insertChatToDB(tx *sql.Tx, chat *Chat) error {
	err := tx.QueryRow(`INSERT INTO chats (name) VALUES ($1) RETURNING id`, chat.Name).Scan(&chat.ID)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update chat in db
func updateChatInDB(tx *sql.Tx, chat *Chat) error {
	_, err := tx.Exec(`UPDATE chats SET name = $1, is_deleted = $2 WHERE id = $3`, chat.Name, chat.IsDeleted, chat.ID)
	if err != nil {
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
		return nil, appErr.InternalServerError("internal server error")
	}
	chat.IsGroupChat = false

	return &chat, nil
}