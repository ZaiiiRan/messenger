package chatMember

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/shortUser"
	"database/sql"
	"fmt"
	"time"
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
		logger.GetInstance().Error(err.Error(), "get chat member role from db", map[string]interface{}{"memberID": memberID, "chatID": chatID}, err)
		return "not member", appErr.InternalServerError("internal server error")
	}
	return role, nil
}

// get chat member fields: removed_by, removed_at, added_by and removed_at
func getChatMemberRemoveAndAddInfo(memberID uint64, chatID uint64) (*uint64, *time.Time, uint64, *time.Time, error) {
	db := pgDB.GetDB()
	var removedBy *uint64
	var addedBy uint64
	var addedAt time.Time
	var removedAt *time.Time

	err := db.QueryRow(`SELECT removed_by, added_by, added_at, removed_at FROM chat_members WHERE chat_id = $1 AND user_id = $2`, chatID, memberID).Scan(&removedBy,
		&addedBy, &addedAt, &removedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil, 0, nil, appErr.NotFound(fmt.Sprintf("user with id %d in chat with id %d not found", memberID, chatID))
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get chat member remove and add info", map[string]interface{}{"memberID": memberID, "chatID": chatID}, err)
		return nil, nil, 0, nil, appErr.InternalServerError("internal server error")
	}

	return removedBy, removedAt, addedBy, &addedAt, nil
}

// get role id from db
func getRoleIDFromDB(roleString string) (int, error) {
	db := pgDB.GetDB()
	var roleID int
	err := db.QueryRow(`SELECT id FROM chat_roles WHERE role = $1`, roleString).Scan(&roleID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "get chat role id by role string", map[string]interface{}{"roleString": roleString}, err)
		return 0, appErr.InternalServerError("internal server error")
	}
	return roleID, nil
}

// insert chat member to db
func insertChatMemberToDB(tx *sql.Tx, member *ChatMember, roleID int) error {
	_, err := tx.Exec(`INSERT INTO chat_members (chat_id, user_id, role_id, added_by, added_at) VALUES ($1, $2, $3, $4, $5)`,
		member.ChatID, member.User.ID, roleID, member.AddedBy, member.AddedAt)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "insert chat member to db", member, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update chat member in db
func updateChatMemberInDB(tx *sql.Tx, member *ChatMember, roleID int) error {
	_, err := tx.Exec(`UPDATE chat_members SET role_id = $1, removed_by = $2, added_by = $3, added_at = $4, removed_at = $5 WHERE chat_id = $6 AND user_id = $7`,
		roleID, member.RemovedBy, member.AddedBy, member.AddedAt, member.RemovedAt, member.ChatID, member.User.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "update chat member in db", member, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// get chat members count from db
func getChatMembersCountFromDB(chatID uint64) (int, error) {
	db := pgDB.GetDB()
	var count int
	err := db.QueryRow(`SELECT COUNT(*) as count FROM chat_members WHERE chat_id = $1 AND (removed_by IS NULL OR removed_by = user_id)`, chatID).Scan(&count)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "get chat members count from db", chatID, err)
		return 0, appErr.InternalServerError("internal server error")
	}
	return count, err
}

// Get chat members from db by search string
func getChatMembersFromDB(actorID, chatID uint64) ([]ChatMember, error) {
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname,
			u.is_deleted, u.is_banned, u.is_activated,
			cm.chat_id, cr.role, cm.removed_by, cm.added_by, cm.removed_at
		FROM chat_members cm
		JOIN users u ON u.id = cm.user_id
		JOIN chat_roles cr ON cr.id = cm.role_id
		WHERE 
			cm.chat_id = $1
			AND cm.user_id != $2
			AND (cm.removed_by IS NULL OR cm.removed_by = cm.user_id)
		ORDER BY CASE 
			WHEN cm.removed_by = cm.user_id THEN 1
			ELSE 0
		END, added_at
	`

	return queryMembers(query, chatID, actorID)
}

// search members query executing
func queryMembers(query string, params ...interface{}) ([]ChatMember, error) {
	db := pgDB.GetDB()

	rows, err := db.Query(query, params...)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("members not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "query members", query, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	defer rows.Close()

	members, err := createMembersFromSQLRows(rows)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, appErr.NotFound("members not found")
	}

	return members, nil
}

// creating members from sql rows
func createMembersFromSQLRows(rows *sql.Rows) ([]ChatMember, error) {
	var members []ChatMember

	for rows.Next() {
		var member ChatMember
		var shortUser shortUser.ShortUser
		var roleString *string
		err := rows.Scan(&shortUser.ID, &shortUser.Username, &shortUser.Firstname, &shortUser.Lastname,
			&shortUser.IsDeleted, &shortUser.IsBanned, &shortUser.IsActivated,
			&member.ChatID, &roleString, &member.RemovedBy, &member.AddedBy, &member.RemovedAt)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "creating chat members from sql rows", rows, err)
			return nil, appErr.InternalServerError("internal server error")
		}
		member.User = &shortUser
		member.Role = GetRoleValue(roleString)

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		logger.GetInstance().Error(err.Error(), "creating chat members from sql rows", rows, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return members, nil
}
