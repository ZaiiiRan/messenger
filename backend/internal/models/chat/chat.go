package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/models/shortUser"
	"backend/internal/models/socialUser"
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
	User      *shortUser.ShortUser `json:"user"`
	Role      int                  `json:"role"`
	RemovedBy *uint64
	AddedBy   uint64
	ChatID    uint64 `json:"chat_id"`
}

// Creating chat object with validations (for inserting)
func CreateChat(name string, members []uint64, isGroup bool, ownerDTO *user.UserDTO) (*Chat, []ChatMember, error) {
	var chatName *string
	if isGroup {
		if err := validateName(name); err != nil {
			return nil, nil, err
		}
		chatName = &name
		if len(members) < 2 {
			return nil, nil, appErr.BadRequest("need at least 2 members for group chat")
		}
	} else if len(members) < 1 {
		return nil, nil, appErr.BadRequest("need at least 1 member for private chat")
	}

	chat := &Chat{
		Name:        chatName,
		IsGroupChat: isGroup,
	}

	var chatMembers []ChatMember
	ownerRole := Roles.Owner
	if !isGroup {
		ownerRole = Roles.Member
	}

	chatMembers = append(chatMembers, ChatMember{
		User:    shortUser.CreateShortUserFromUserDTO(ownerDTO),
		Role:    ownerRole,
		AddedBy: ownerDTO.ID,
	})

	for _, memberID := range members {
		user, err := user.GetUserByID(memberID)
		if err != nil {
			if err.Error() == "user not found" {
				return nil, nil, appErr.BadRequest(fmt.Sprintf("user with id %d not found", memberID))
			}
			return nil, nil, err
		}
		relations, err := socialUser.GetRelations(ownerDTO.ID, memberID)
		if err != nil {
			return nil, nil, err
		}
		if (relations != nil && *relations != "accepted") || (relations == nil) {
			return nil, nil, appErr.BadRequest(fmt.Sprintf("user with id %d is not in your friends list", memberID))
		}
		chatMembers = append(chatMembers, ChatMember{
			User:    shortUser.CreateShortUserFromUser(user),
			Role:    Roles.Member,
			AddedBy: ownerDTO.ID,
		})
	}

	return chat, chatMembers, nil
}

// Save chat with members
func (chat *Chat) SaveWithMembers(newMembers []ChatMember) ([]ChatMember, error) {
	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	if err := chat.saveChatToDB(tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	var members []ChatMember
	for _, member := range newMembers {
		if err := chat.saveMemberToDB(tx, &member); err != nil {
			tx.Rollback()
			return nil, err
		}
		member.ChatID = chat.ID
		members = append(members, member)
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return members, nil
}

func (chat *Chat) Rename(newName string, sender uint64) error {
	if !chat.IsGroupChat {
		return appErr.BadRequest("chat is not a group chat")
	}

	self, err := chat.GetChatMemberByID(sender)
	if err != nil {
		return err
	}
	if self.RemovedBy != nil || self.Role < Roles.Admin {
		return appErr.BadRequest("you don't have enough rights")
	}

	if *chat.Name == newName {
		return appErr.BadRequest("the names are the same")
	}

	err = validateName(newName)
	if err != nil {
		return err
	}

	chat.Name = &newName

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
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

	if err := tx.Commit(); err != nil {
		return appErr.InternalServerError("internal server error")
	}

	return nil
}

// Return user to chat
func (chat *Chat) ReturnToChat(memberID uint64) (*ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}
	self, err := chat.GetChatMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	if self.RemovedBy == nil {
		return nil, appErr.BadRequest("you are already in chat")
	}
	if *self.RemovedBy != memberID {
		return nil, appErr.BadRequest("you have been excluded from the chat")
	}

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	self.RemovedBy = nil
	if self.Role != Roles.Owner {
		self.Role = Roles.Member
	}

	err = chat.saveMemberToDB(tx, self)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return self, nil
}

// Add members to chat
func (chat *Chat) AddMembers(newMembersIDs []uint64, senderID uint64) ([]ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}
	self, err := chat.GetChatMemberByID(senderID)
	if err != nil {
		return nil, err
	}
	if self.RemovedBy != nil {
		return nil, appErr.BadRequest("you don't have enough rights")
	}

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	var added []ChatMember
	for _, memberID := range newMembersIDs {
		member, err := chat.addMember(tx, memberID, senderID)
		if err != nil {
			return nil, err
		}
		added = append(added, *member)
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return added, nil
}

// add 1 member to chat
func (chat *Chat) addMember(tx *sql.Tx, targetID, senderID uint64) (*ChatMember, error) {
	user, err := user.GetUserByID(targetID)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, appErr.BadRequest(fmt.Sprintf("user with id %d not found", targetID))
		}
		return nil, err
	}

	relations, err := socialUser.GetRelations(senderID, targetID)
	if err != nil {
		return nil, err
	}

	if (relations != nil && *relations != "accepted") || (relations == nil) {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d is not in your friends list", targetID))
	}

	target, err := chat.GetChatMemberByID(targetID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if target != nil && target.RemovedBy != nil {
		target.RemovedBy = nil
		target.Role = Roles.Member
		target.AddedBy = senderID

		err = chat.saveMemberToDB(tx, target)
		if err != nil {
			return nil, err
		}
		return target, nil
	}

	if target != nil && target.RemovedBy == nil {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d is already a chat member", targetID))
	}

	newMember := &ChatMember{
		User:    shortUser.CreateShortUserFromUser(user),
		AddedBy: senderID,
		Role:    Roles.Member,
	}

	err = chat.saveMemberToDB(tx, newMember)
	if err != nil {
		return nil, err
	}

	return newMember, nil
}

// Remove members from chat
func (chat *Chat) RemoveMembers(membersIDs []uint64, senderID uint64) ([]ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}
	self, err := chat.GetChatMemberByID(senderID)
	if err != nil {
		return nil, err
	}
	if self.RemovedBy != nil {
		return nil, appErr.BadRequest("you don't have enough rights")
	}

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	var removed []ChatMember
	for _, memberID := range membersIDs {
		member, err := chat.removeMember(tx, memberID, senderID)
		if err != nil {
			return nil, err
		}
		removed = append(removed, *member)
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return removed, nil
}

// remove 1 member from chat
func (chat *Chat) removeMember(tx *sql.Tx, memberID, senderID uint64) (*ChatMember, error) {
	member, err := chat.GetChatMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	if member.RemovedBy != nil && member.User.ID != *member.RemovedBy {
		return nil, appErr.BadRequest(fmt.Sprintf("the member with id %d has already been deleted", memberID))
	} else if member.RemovedBy != nil && member.User.ID == *member.RemovedBy && *member.RemovedBy != senderID {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d has left the chat", memberID))
	} else if member.RemovedBy != nil && member.User.ID == *member.RemovedBy && *member.RemovedBy == senderID {
		return nil, appErr.BadRequest("you have already left the chat")
	}

	if member.AddedBy != senderID {
		removerRole, err := chat.GetMemberRole(senderID)
		if err != nil {
			return nil, err
		}
		if removerRole == Roles.NotMember {
			return nil, appErr.InternalServerError("internal server error")
		}

		if removerRole < member.Role {
			return nil, appErr.BadRequest("you don't have enough rights")
		}
	} else if (member.AddedBy == senderID && member.User.ID != senderID) && member.Role > Roles.Member {
		return nil, appErr.BadRequest("you don't have enough rights")
	}
	member.RemovedBy = &senderID
	if member.Role != Roles.Owner {
		member.Role = Roles.Member
	}

	err = chat.saveMemberToDB(tx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// Get chat member role
func (chat *Chat) GetMemberRole(memberID uint64) (int, error) {
	db := pgDB.GetDB()
	var role *string
	err := db.QueryRow(`
		SELECT cr.role
		FROM chat_members cm
		LEFT JOIN chat_roles cr ON cm.role_id = cr.id
		WHERE cm.user_id = $1`, memberID).Scan(&role)
	if err != nil && err == sql.ErrNoRows {
		role = nil
	} else if err != nil {
		fmt.Println(err)
		return Roles.NotMember, appErr.InternalServerError("internal server error")
	}
	roleValue := GetRoleValue(role)
	return roleValue, nil
}

// Get chat member by id from db
func (chat *Chat) GetChatMemberByID(targetID uint64) (*ChatMember, error) {
	db := pgDB.GetDB()
	var member ChatMember
	role, err := chat.GetMemberRole(targetID)
	if err != nil {
		return nil, err
	}
	if role == Roles.NotMember {
		return nil, appErr.NotFound(fmt.Sprintf("user with id %d in chat with id %d not found", targetID, chat.ID))
	}
	member.Role = role

	user, err := user.GetUserByID(targetID)
	if err != nil && err.Error() == "user not found" {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d not found", targetID))
	} else if err != nil {
		return nil, err
	}
	shortUser := shortUser.CreateShortUserFromUser(user)
	member.User = shortUser
	member.ChatID = chat.ID

	err = db.QueryRow(`SELECT removed_by, added_by FROM chat_members WHERE chat_id = $1 AND user_id = $2`, chat.ID, targetID).Scan(&member.RemovedBy, &member.AddedBy)
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d in chat with id %d not found", targetID, chat.ID))
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	return &member, nil
}

func (chat *Chat) GetAllChatMembers() ([]ChatMember, error) {
	db := pgDB.GetDB()
	rows, err := db.Query(`
		SELECT cm.user_id, cr.role, u.username, u.firstname, u.lastname, u.is_deleted, u.is_banned, u.is_activated
		FROM chat_members cm
		JOIN chat_roles cr ON cm.role_id = cr.id
		JOIN users u ON cm.user_id = u.id
		WHERE cm.chat_id = $1`, chat.ID)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("chat members not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	defer rows.Close()

	var members []ChatMember
	for rows.Next() {
		var member ChatMember
		var user shortUser.ShortUser

		err := rows.Scan(
			&user.ID,
			&member.Role,
			&user.Username,
			&user.Firstname,
			&user.Lastname,
			&user.IsDeleted,
			&user.IsBanned,
			&user.IsActivated,
		)
		if err != nil {
			return nil, appErr.InternalServerError("internal server error")
		}

		member.User = &user
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return members, nil
}

// Get chat by id from db
func GetChatByID(id uint64) (*Chat, error) {
	db := pgDB.GetDB()

	var chat Chat
	err := db.QueryRow(`SELECT id, name, is_deleted FROM chats WHERE id = $1 AND is_deleted != TRUE`, id).Scan(&chat.ID, &chat.Name, &chat.IsDeleted)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("chat not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	if chat.Name != nil && *chat.Name != "" {
		chat.IsGroupChat = true
	}

	return &chat, nil
}

// validate chat name
func validateName(name string) error {
	if name == "" {
		return appErr.BadRequest("chat name is empty")
	}
	if len(name) < 5 {
		return appErr.BadRequest("chat name must be at least 5 characters")
	}
	return nil
}

// save chat in db
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

// save chat member in db
func (chat *Chat) saveMemberToDB(tx *sql.Tx, member *ChatMember) error {
	isInserting := member.ChatID == 0

	roleString := GetRoleString(member.Role)
	if roleString == "" {
		return appErr.InternalServerError("internal server error")
	}
	var roleID int
	err := tx.QueryRow(`SELECT id FROM chat_roles WHERE role = $1`, roleString).Scan(&roleID)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}

	if isInserting {
		_, err = tx.Exec(`INSERT INTO chat_members (chat_id, user_id, role_id, added_by) VALUES ($1, $2, $3, $4)`, chat.ID, member.User.ID, roleID, member.AddedBy)
		if err != nil {
			return appErr.InternalServerError("internal server error")
		}
	} else {
		_, err := tx.Exec(`UPDATE chat_members SET role_id = $1, removed_by = $2, added_by = $3 WHERE chat_id = $4 AND user_id = $5`,
			roleID, member.RemovedBy, member.AddedBy, chat.ID, member.User.ID)
		if err != nil {
			return appErr.InternalServerError("internal server error")
		}
	}
	return nil
}
