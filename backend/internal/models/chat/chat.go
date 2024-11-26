package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/models/shortUser"
	"backend/internal/models/user"
	"database/sql"
	"fmt"
)

type Chat struct {
	ID          uint64  `json:"id"`
	Name        *string `json:"name"`
	IsGroupChat bool    `json:"is_group_chat"`
	IsDeleted   bool    `json:"is_deleted"`
}

type ChatMember struct {
	User *shortUser.ShortUser `json:"user"`
	Role string               `json:"role"`
}

func CreateChat(name string, members []uint64, isGroup bool, ownerDTO *user.UserDTO) (*Chat, []ChatMember, error) {
	var chatName *string
	if !isGroup {
		if len(members) < 1 {
			return nil, nil, appErr.BadRequest("need at least 1 member")
		}
	} else {
		err := validateName(name)
		if err != nil {
			return nil, nil, err
		}
		chatName = &name
		if len(members) < 2 {
			return nil, nil, appErr.BadRequest("need at least 2 members")
		}
	}

	var users []ChatMember

	ownerRole := "owner"
	if !isGroup {
		ownerRole = "member"
	}
	users = append(users, ChatMember{
		User: shortUser.CreateShortUserFromUserDTO(ownerDTO),
		Role: ownerRole,
	})

	for _, memberID := range members {
		user, err := user.GetUserByID(memberID)
		if err != nil {
			return nil, nil, err
		}
		member := ChatMember{
			User: shortUser.CreateShortUserFromUser(user),
			Role: "member",
		}
		users = append(users, member)
	}

	return &Chat{
		Name:        chatName,
		IsGroupChat: isGroup,
	}, users, nil
}

func (chat *Chat) Save(members []ChatMember) error {
	isCreating := chat.ID == 0

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		fmt.Println(err)
		return appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()
	err = chat.saveChatToDB(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if isCreating {
		for _, member := range members {
			err = chat.addMemberToDB(tx, &member)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return appErr.InternalServerError("internal server error")
	}

	return nil
}

func (chat *Chat) AddMember(targetID uint64) error {
	// if !chat.IsGroupChat {
	// 	return appErr.BadRequest("chat is not a group chat")
	// }
	// if chat.isMember(targetID) {
	// 	return appErr.BadRequest(fmt.Sprintf("user with id %d is already a member", targetID))
	// }
	// user, err := user.GetUserByID(targetID)
	// if err != nil {
	// 	return err
	// }
	// member := ChatMember{
	// 	User: shortUser.CreateShortUserFromUser(user),
	// 	Role: "member",
	// }
	// chat.Members = append(chat.Members, member)
	return nil
}

func (chat *Chat) isMember(targetID uint64) bool {
	// for _, member := range chat.Members {
	// 	if member.User.ID == targetID {
	// 		return true
	// 	}
	// }
	return false
}

func GetChatByID(id uint64) (*Chat, error) {
	db := pgDB.GetDB()

	var chat Chat
	err := db.QueryRow(`SELECT id, name, is_deleted FROM chats WHERE id = $1`, id).Scan(&chat.ID, &chat.Name, &chat.IsDeleted)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("chat not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	if chat.Name != nil && *chat.Name != "" {
		chat.IsGroupChat = true
	}

	// rows, err := db.Query(`
	// 	SELECT cm.user_id, cr.role, u.username, u.firstname, u.lastname, u.is_deleted, u.is_banned, u.is_activated
	// 	FROM chat_members cm
	// 	JOIN chat_roles cr ON cm.role_id = cr.id
	// 	JOIN users u ON cm.user_id = u.id
	// 	WHERE cm.chat_id = $1`, id)
	// if err != nil {
	// 	return nil, appErr.InternalServerError("internal server error")
	// }

	// defer rows.Close()

	// var members []ChatMember
	// for rows.Next() {
	// 	var member ChatMember
	// 	var user shortUser.ShortUser

	// 	err := rows.Scan(
	// 		&user.ID,
	// 		&member.Role,
	// 		&user.Username,
	// 		&user.Firstname,
	// 		&user.Lastname,
	// 		&user.IsDeleted,
	// 		&user.IsBanned,
	// 		&user.IsActivated,
	// 	)
	// 	if err != nil {
	// 		return nil, appErr.InternalServerError("internal server error")
	// 	}

	// 	member.User = &user
	// 	members = append(members, member)
	// }
	// if err = rows.Err(); err != nil {
	// 	return nil, appErr.InternalServerError("internal server error")
	// }

	// chat.Members = members

	return &chat, nil
}

func validateName(name string) error {
	if name == "" {
		return appErr.BadRequest("chat name is empty")
	}
	if len(name) < 5 {
		return appErr.BadRequest("chat name must be at least 5 characters")
	}
	return nil
}

func (chat *Chat) saveChatToDB(tx *sql.Tx) error {
	if chat.ID == 0 {
		err := tx.QueryRow(`INSERT INTO chats (name) VALUES ($1) RETURNING id`, chat.Name).Scan(&chat.ID)
		if err != nil {
			return appErr.InternalServerError("internal server error")
		}
	} else {
		_, err := tx.Exec(`UPDATE chats SET name = $1, is_deleted = $2 WHERE id = $3`, chat.Name, chat.IsDeleted, chat.ID)
		if err != nil {
			return appErr.InternalServerError("internal server error")
		}
	}
	return nil
}

func (chat *Chat) addMemberToDB(tx *sql.Tx, member *ChatMember) error {
	var roleID int
	err := tx.QueryRow(`SELECT id FROM chat_roles WHERE role = $1`, member.Role).Scan(&roleID)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}
	_, err = tx.Exec(`INSERT INTO chat_members (chat_id, user_id, role_id) VALUES ($1, $2, $3)`, chat.ID, member.User.ID, roleID)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}
	return nil
}
