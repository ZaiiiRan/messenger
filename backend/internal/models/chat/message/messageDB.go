package message

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
)

// insert message to db
func insertMessageToDB(message *Message) error {
	db := pgDB.GetDB()
	query := `
		INSERT INTO messages (chat_id, user_id, content) VALUES
		($1, $2, $3)
		RETURNING id, sent_at
	`
	err := db.QueryRow(query, message.Chat.ID, message.ChatMember.User.ID, message.Content).Scan(&message.ID, &message.SentAt)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "insert message to db", message, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update message in db
func updateMessageInDB(message *Message) error {
	db := pgDB.GetDB()
	query := `
		UPDATE messages
		SET
			content = $1,
			last_edited = $2,
			is_deleted = $3
		WHERE id = $4
	`
	_, err := db.Exec(query, message.Content, message.LastEdited, message.IsDeleted, message.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "update message in db", message, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}
