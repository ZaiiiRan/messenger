package chatMember

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"database/sql"
	"fmt"
)

// get chat member role from db by member id and chat id
func getChatMemberRoleFromDB(memberID uint64, chatID uint64) (string, error) {
	db := pgDB.GetDB()
	var role string
	err := db.QueryRow(`
		SELECT cr.role
		FROM chat_members cm
		LEFT JOIN chat_roles cr ON cm.role_id = cr.id
		WHERE cm.user_id = $1 AND cm.chat_id = $2`, memberID, chatID).Scan(&role)
	if err != nil && err == sql.ErrNoRows {
		role = "not member"
	} else if err != nil {
		return "not member", appErr.InternalServerError("internal server error")
	}
	return role, nil
}

// get chat member fields: removed_by and added_by
func getChatMemberRemoveAndAddInfo(memberID uint64, chatID uint64) (*uint64, uint64, error) {
	db := pgDB.GetDB()
	var removedBy *uint64
	var addedBy uint64

	err := db.QueryRow(`SELECT removed_by, added_by FROM chat_members WHERE chat_id = $1 AND user_id = $2`, chatID, memberID).Scan(&removedBy, &addedBy)
	if err != nil && err == sql.ErrNoRows {
		return nil, 0, appErr.NotFound(fmt.Sprintf("user with id %d in chat with id %d not found", memberID, chatID))
	} else if err != nil {
		return nil, 0, appErr.InternalServerError("internal server error")
	}

	return removedBy, addedBy, nil
}

// get role id from db
func getRoleIDFromDB(roleString string) (int, error) {
	db := pgDB.GetDB()
	var roleID int
	err := db.QueryRow(`SELECT id FROM chat_roles WHERE role = $1`, roleString).Scan(&roleID)
	if err != nil {
		return 0, appErr.InternalServerError("internal server error")
	}
	return roleID, nil
}

// insert chat member to db
func insertChatMemberToDB(tx *sql.Tx, member *ChatMember, roleID int) error {
	_, err := tx.Exec(`INSERT INTO chat_members (chat_id, user_id, role_id, added_by) VALUES ($1, $2, $3, $4)`, member.ChatID, member.User.ID, roleID, member.AddedBy)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update chat member in db
func updateChatMemberInDB(tx *sql.Tx, member *ChatMember, roleID int) error {
	_, err := tx.Exec(`UPDATE chat_members SET role_id = $1, removed_by = $2, added_by = $3 WHERE chat_id = $4 AND user_id = $5`,
		roleID, member.RemovedBy, member.AddedBy, member.ChatID, member.User.ID)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}
	return nil
}
