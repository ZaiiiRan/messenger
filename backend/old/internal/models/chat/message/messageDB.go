package message

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/shortUser"
	"backend/internal/models/user"
	"database/sql"
	"errors"
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

// get message by id from db
func getMessageByIDFromDB(id uint64) (*Message, error) {
	db := pgDB.GetDB()
	query := `
		SELECT m.id, chat_id, m.user_id, m.content, m.sent_at, m.last_edited, m.is_deleted
		FROM messages m
		WHERE m.id = $1 AND m.is_deleted = FALSE
	`

	row := db.QueryRow(query, id)

	return createMessageFromSQLRow(row)
}

// get messages from db
func getMessagesFromDB(chat *chat.Chat, actor *chatMember.ChatMember, limit, offset int) ([]Message, error) {
	query := `
		WITH membership_period AS (
			SELECT 
				cm.added_at, 
				COALESCE(cm.removed_at, NOW()) AS removed_at
			FROM chat_members cm
			WHERE 
				cm.user_id = $2 
				AND cm.chat_id = $1
		)
		SELECT 
			m.id, m.user_id, m.content, m.sent_at, m.last_edited 
		FROM messages m
		JOIN membership_period mp ON m.sent_at BETWEEN mp.added_at AND mp.removed_at
		WHERE
			m.chat_id = $1
			AND m.is_deleted = FALSE
		ORDER BY m.sent_at DESC
		LIMIT $3 OFFSET $4
	`

	return queryMessages(chat, query, chat.ID, actor.User.ID, limit, offset)
}

// query messages
func queryMessages(chat *chat.Chat, query string, params ...interface{}) ([]Message, error) {
	db := pgDB.GetDB()

	rows, err := db.Query(query, params...)
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.NotFound("messages not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "query messages", query, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	defer rows.Close()

	messages, err := createMessagesFromSQLRows(chat, rows)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, appErr.NotFound("messages not found")
	}

	return messages, nil
}

// parsing messages from sql rows
func createMessagesFromSQLRows(chat *chat.Chat, rows *sql.Rows) ([]Message, error) {
	var messages []Message

	for rows.Next() {
		var message Message
		var memberID uint64
		err := rows.Scan(&message.ID, &memberID, &message.Content,
			&message.SentAt, &message.LastEdited)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "creating messages from sql rows", rows, err)
			return nil, appErr.InternalServerError("internal server error")
		}
		message.Chat = chat

		var appError *appErr.AppError
		member, err := chat.GetChatMemberByID(memberID)
		if err != nil && (errors.As(err, &appError) && appError.StatusCode == 404) {
			member = &chatMember.ChatMember{}
			member.Role = chatMember.Roles.NotMember
			user, err := user.GetUserByID(memberID)
			if err != nil {
				return nil, err
			}
			memberShortUser := shortUser.CreateShortUserFromUser(user)
			member.User = memberShortUser
		} else if err != nil {
			return nil, err
		}

		message.ChatMember = member

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		logger.GetInstance().Error(err.Error(), "creating messages from sql rows", rows, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return messages, nil
}

// create message from sql row
func createMessageFromSQLRow(row *sql.Row) (*Message, error) {
	var message Message
	var chatID uint64
	var memberID uint64
	err := row.Scan(&message.ID, &chatID, &memberID, &message.Content,
		&message.SentAt, &message.LastEdited, &message.IsDeleted)
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.NotFound("message not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "creating message from sql row", row, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	chat, err := chat.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}
	message.Chat = chat

	member, err := chat.GetChatMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	message.ChatMember = member

	return &message, nil
}
