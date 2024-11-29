package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/models/chatMember"
	"backend/internal/models/shortUser"
	"backend/internal/models/socialUser"
	"backend/internal/models/user"
	"database/sql"
	"errors"
	"fmt"
)

type Chat struct {
	ID          uint64  `json:"id"`
	Name        *string `json:"name"`
	IsGroupChat bool    `json:"is_group_chat"`
	IsDeleted   bool    `json:"is_deleted"`
}

// Creating chat object with validations (for inserting)
func CreateChat(name string, members []uint64, isGroup bool, ownerDTO *user.UserDTO) (*Chat, []chatMember.ChatMember, error) {
	err := validateBeforeCreatingChat(name, members, isGroup, ownerDTO)
	if err != nil {
		return nil, nil, err
	}

	var chatName *string
	if isGroup {
		chatName = &name
	}

	chat := &Chat{Name: chatName, IsGroupChat: isGroup}

	var chatMembers []chatMember.ChatMember
	ownerRole := chatMember.Roles.Owner
	if !isGroup {
		ownerRole = chatMember.Roles.Member
	}

	chatMembers = append(chatMembers, chatMember.ChatMember{
		User:    shortUser.CreateShortUserFromUserDTO(ownerDTO),
		Role:    ownerRole,
		AddedBy: ownerDTO.ID,
	})

	for _, memberID := range members {
		user, err := getUserForAdding(memberID, ownerDTO.ID)
		if err != nil {
			return nil, nil, err
		}

		chatMembers = append(chatMembers, chatMember.ChatMember{
			User:    shortUser.CreateShortUserFromUser(user),
			Role:    chatMember.Roles.Member,
			AddedBy: ownerDTO.ID,
		})
	}

	return chat, chatMembers, nil
}

// Save chat with members
func (chat *Chat) SaveWithMembers(newMembers []chatMember.ChatMember) ([]chatMember.ChatMember, error) {
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

	var members []chatMember.ChatMember
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

// Rename chat
func (chat *Chat) Rename(newName string, actor *chatMember.ChatMember) error {
	if !chat.IsGroupChat {
		return appErr.BadRequest("chat is not a group chat")
	}

	if actor.Role < chatMember.Roles.Admin {
		return appErr.Forbidden("you don't have enough rights")
	}

	if *chat.Name == newName {
		return appErr.BadRequest("the names are the same")
	}

	err := validateName(newName)
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

func (chat *Chat) Delete(actor *chatMember.ChatMember) error {
	if !chat.IsGroupChat {
		return appErr.BadRequest("chat is not a group chat")
	}

	if actor.Role < chatMember.Roles.Owner {
		return appErr.Forbidden("you don't have enough rights")
	}

	chat.IsDeleted = true

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

// Chat member role changing
func (chat *Chat) ChatMemberRoleChange(memberID uint64, newRole string, actor *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}

	if actor.Role == chatMember.Roles.Member {
		return nil, appErr.Forbidden("you don't have enough rights")
	}
	if memberID == actor.User.ID {
		return nil, appErr.Forbidden("you can't change your role")
	}

	member, err := chat.GetChatMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	if member.Removed() {
		return nil, appErr.BadRequest("you cannot assign a role to an excluded member")
	}

	roleValue := chatMember.GetRoleValue(&newRole)
	if roleValue == chatMember.Roles.NotMember {
		return nil, appErr.BadRequest("unknown role")
	}
	if roleValue == chatMember.Roles.Owner {
		return nil, appErr.BadRequest("owner role cannot be assigned")
	}

	if actor.Role <= roleValue || (member.Role == actor.Role && roleValue < member.Role) || member.Role > actor.Role {
		return nil, appErr.Forbidden("you don't have enough rights")
	}
	if member.Role == roleValue {
		return nil, appErr.BadRequest(fmt.Sprintf("member with id %d is already %s", memberID, newRole))
	}

	member.Role = roleValue

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	err = chat.saveMemberToDB(tx, member)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return member, nil
}

// Leave from chat
func (chat *Chat) LeaveFromChat(member *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}

	if member.RemovedBy != nil && *member.RemovedBy != member.User.ID {
		return nil, appErr.Forbidden("you have been removed from the chat")
	}

	member.RemovedBy = &member.User.ID
	if member.Role != chatMember.Roles.Owner {
		member.Role = chatMember.Roles.Member
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

	err = chat.saveMemberToDB(tx, member)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return member, nil
}

// Return a left user to chat
func (chat *Chat) ReturnToChat(member *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}

	if member.RemovedBy == nil {
		return nil, appErr.BadRequest("you are already in chat")
	}
	if *member.RemovedBy != member.User.ID {
		return nil, appErr.Forbidden("you have been removed from the chat")
	}

	member.RemovedBy = nil
	if member.Role != chatMember.Roles.Owner {
		member.Role = chatMember.Roles.Member
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

	err = chat.saveMemberToDB(tx, member)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return member, nil
}

// Add members to chat
func (chat *Chat) AddMembers(newMembersIDs []uint64, addingMember *chatMember.ChatMember) ([]chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
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

	var added []chatMember.ChatMember
	for _, memberID := range newMembersIDs {
		if addingMember.User.ID == memberID {
			tx.Rollback()
			return nil, appErr.BadRequest("you can't add yourself")
		}

		member, err := chat.addMember(tx, memberID, addingMember)
		if err != nil {
			tx.Rollback()
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
func (chat *Chat) addMember(tx *sql.Tx, targetID uint64, addingMember *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	user, err := getUserForAdding(targetID, addingMember.User.ID)
	if err != nil {
		return nil, err
	}

	var appError *appErr.AppError
	target, err := chat.GetChatMemberByID(targetID)
	if err != nil && (errors.As(err, &appError) && appError.StatusCode != 404) {
		return nil, err
	}

	if target != nil && target.Removed() {
		target, err = chat.addOldMemberToChat(tx, target, addingMember)
		if err != nil {
			return nil, err
		}
		return target, nil
	}

	if target != nil && !target.Removed() {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d is already a chat member", targetID))
	}

	target, err = chat.addNewMemberToChat(tx, user, addingMember)
	if err != nil {
		return nil, err
	}
	return target, err
}

// returning old member to chat
func (chat *Chat) addOldMemberToChat(tx *sql.Tx, target, addingMember *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	target.RemovedBy = nil
	target.Role = chatMember.Roles.Member
	target.AddedBy = addingMember.User.ID

	err := chat.saveMemberToDB(tx, target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// adding new member to chat
func (chat *Chat) addNewMemberToChat(tx *sql.Tx, user *user.User, addingMember *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	newMember := &chatMember.ChatMember{
		User:    shortUser.CreateShortUserFromUser(user),
		AddedBy: addingMember.User.ID,
		Role:    chatMember.Roles.Member,
	}
	err := chat.saveMemberToDB(tx, newMember)
	if err != nil {
		return nil, err
	}

	return newMember, nil
}

// Remove members from chat
func (chat *Chat) RemoveMembers(membersIDs []uint64, removingMember *chatMember.ChatMember) ([]chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
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

	var removed []chatMember.ChatMember
	for _, memberID := range membersIDs {
		if removingMember.User.ID == memberID {
			tx.Rollback()
			return nil, appErr.BadRequest("you can't remove yourself")
		}

		member, err := chat.removeMember(tx, memberID, removingMember)
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
func (chat *Chat) removeMember(tx *sql.Tx, memberID uint64, removingMember *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	member, err := chat.GetChatMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	if member.RemovedBy != nil && member.User.ID != *member.RemovedBy {
		return nil, appErr.BadRequest(fmt.Sprintf("the member with id %d has already been deleted", memberID))
	} else if member.RemovedBy != nil && member.User.ID == *member.RemovedBy {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d has left the chat", memberID))
	}

	if member.AddedBy != removingMember.User.ID {
		if removingMember.Role == chatMember.Roles.NotMember {
			return nil, appErr.InternalServerError("internal server error")
		}

		if removingMember.Role <= member.Role {
			return nil, appErr.Forbidden("trying to delete a member with a higher role")
		}
	} else if member.AddedBy == removingMember.User.ID && member.Role >= removingMember.Role {
		return nil, appErr.BadRequest("trying to delete a member with a higher role")
	}
	member.RemovedBy = &removingMember.User.ID

	err = chat.saveMemberToDB(tx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// Get chat member role
func (chat *Chat) GetMemberRole(memberID uint64) (int, error) {
	db := pgDB.GetDB()
	var role string
	err := db.QueryRow(`
		SELECT cr.role
		FROM chat_members cm
		LEFT JOIN chat_roles cr ON cm.role_id = cr.id
		WHERE cm.user_id = $1 AND cm.chat_id = $2`, memberID, chat.ID).Scan(&role)
	if err != nil && err == sql.ErrNoRows {
		role = "not member"
	} else if err != nil {
		return chatMember.Roles.NotMember, appErr.InternalServerError("internal server error")
	}
	roleValue := chatMember.GetRoleValue(&role)
	return roleValue, nil
}

// Get chat member by id from db
func (chat *Chat) GetChatMemberByID(targetID uint64) (*chatMember.ChatMember, error) {
	db := pgDB.GetDB()
	var member chatMember.ChatMember
	role, err := chat.GetMemberRole(targetID)
	if err != nil {
		return nil, err
	}
	if role == chatMember.Roles.NotMember {
		return nil, appErr.NotFound(fmt.Sprintf("user with id %d in chat with id %d not found", targetID, chat.ID))
	}
	member.Role = role

	user, err := user.GetUserByID(targetID)
	if err != nil && err.Error() == "user not found" {
		return nil, appErr.NotFound(fmt.Sprintf("user with id %d not found", targetID))
	} else if err != nil {
		return nil, err
	}
	shortUser := shortUser.CreateShortUserFromUser(user)
	member.User = shortUser
	member.ChatID = chat.ID

	err = db.QueryRow(`SELECT removed_by, added_by FROM chat_members WHERE chat_id = $1 AND user_id = $2`, chat.ID, targetID).Scan(&member.RemovedBy, &member.AddedBy)
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.NotFound(fmt.Sprintf("user with id %d in chat with id %d not found", targetID, chat.ID))
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	return &member, nil
}

func (chat *Chat) GetAllChatMembers() ([]chatMember.ChatMember, error) {
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

	var members []chatMember.ChatMember
	for rows.Next() {
		var member chatMember.ChatMember
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
	err := db.QueryRow(`SELECT id, name, is_deleted FROM chats WHERE id = $1 AND is_deleted = FALSE`, id).Scan(&chat.ID, &chat.Name, &chat.IsDeleted)
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

// Chat existing checking
func CheckPrivateChatExists(member1, member2 uint64) (bool, error) {
	var chatID uint64
	db := pgDB.GetDB()
	err := db.QueryRow(`
		SELECT c.id FROM chats c
		JOIN chat_members cm ON cm.chat_id = c.id
		WHERE c.name IS NULL 
		AND ((added_by = $1 AND cm.user_id = $2) OR (added_by = $2 AND cm.user_id = $1))
	`, member1, member2).Scan(&chatID)
	if err != nil && err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, appErr.InternalServerError("internal server error")
	}
	return chatID != 0, nil
}

// get user object with access checking
func getUserForAdding(userID uint64, requestSendingMemberID uint64) (*user.User, error) {
	user, err := user.GetUserByID(userID)
	if err != nil {
		var appError *appErr.AppError
		if errors.As(err, &appError) && appError.StatusCode == 404 {
			return nil, appErr.NotFound(fmt.Sprintf("user with id %d not found", userID))
		}
		return nil, err
	}
	err = checkUserAccess(user, requestSendingMemberID)
	if err != nil {
		return nil, err
	}
	return user, nil
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
func (chat *Chat) saveMemberToDB(tx *sql.Tx, member *chatMember.ChatMember) error {
	isInserting := member.ChatID == 0

	roleString := chatMember.GetRoleString(member.Role)
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

// Getting a chat object and its member object (request sender)
func GetChatAndMember(chatID uint64, memberID uint64) (*Chat, *chatMember.ChatMember, error) {
	var appError *appErr.AppError

	chat, err := GetChatByID(chatID)
	if err != nil && errors.As(err, &appError) && appError.StatusCode == 404 {
		return nil, nil, appErr.Forbidden("you cannot access this chat")
	} else if err != nil {
		return nil, nil, err
	}

	member, err := chat.GetChatMemberByID(memberID)
	if err != nil && errors.As(err, &appError) && appError.StatusCode == 404 {
		return nil, nil, appErr.Forbidden("you cannot access this chat")
	} else if err != nil {
		return nil, nil, err
	}

	return chat, member, nil
}

// checking user access (friend status and activation) for adding to chat
func checkUserAccess(target *user.User, requestSendingMemberID uint64) error {
	if target.IsBanned {
		return appErr.BadRequest(fmt.Sprintf("user with id %d is banned", target.ID))
	}
	if !target.IsActivated || target.IsDeleted {
		return appErr.NotFound(fmt.Sprintf("user with id %d not found", target.ID))
	}

	relations, err := socialUser.GetRelations(requestSendingMemberID, target.ID)
	if err != nil {
		return err
	}

	if (relations != nil && *relations != "accepted") || (relations == nil) {
		return appErr.BadRequest(fmt.Sprintf("user with id %d is not in your friends list", target.ID))
	}

	return nil
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

// validate chat name and members count
func validateBeforeCreatingChat(name string, members []uint64, isGroup bool, ownerDTO *user.UserDTO) error {
	if isGroup {
		if err := validateName(name); err != nil {
			return err
		}
		if len(members) < 2 {
			return appErr.BadRequest("need at least 2 members for group chat")
		}
	} else {
		if len(members) < 1 {
			return appErr.BadRequest("need at least 1 member for private chat")
		} else if len(members) > 1 {
			return appErr.BadRequest("max 1 member for private chat")
		}

		privateChatExists, err := CheckPrivateChatExists(ownerDTO.ID, members[0])
		if err != nil {
			return err
		}
		if privateChatExists {
			return appErr.BadRequest(fmt.Sprintf("chat between user with id %d and %d already exists", ownerDTO.ID, members[0]))
		}
	}

	return nil
}
