package chatMember

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/shortUser"
	"backend/internal/models/user"
	"database/sql"
	"fmt"
)

type ChatMember struct {
	User      *shortUser.ShortUser `json:"user"`
	Role      int                  `json:"role"`
	RemovedBy *uint64
	AddedBy   uint64
	ChatID    uint64 `json:"chat_id"`
}

// Removed checking
func (member *ChatMember) Removed() bool {
	return member.RemovedBy != nil
}

// Get chat member role by member id and chat id
func GetChatMemberRole(memberID uint64, chatID uint64) (int, error) {
	roleStr, err := getChatMemberRoleFromDB(memberID, chatID)
	if err != nil {
		return Roles.NotMember, err
	}
	roleValue := GetRoleValue(&roleStr)
	return roleValue, nil
}

// Get chat member by target id and chat id
func GetChatMemberByID(targetID, chatID uint64) (*ChatMember, error) {
	var member ChatMember
	role, err := GetChatMemberRole(targetID, chatID)
	if err != nil {
		return nil, err
	}
	if role == Roles.NotMember {
		return nil, appErr.NotFound(fmt.Sprintf("user with id %d in chat with id %d not found", targetID, chatID))
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
	member.ChatID = chatID

	removedBy, addedBy, err := getChatMemberRemoveAndAddInfo(targetID, chatID)
	if err != nil {
		return nil, err
	}

	member.AddedBy = addedBy
	member.RemovedBy = removedBy

	return &member, nil
}

// Save chat member
func (member *ChatMember) Save(tx *sql.Tx) error {
	isInserting := member.ChatID == 0

	roleString := GetRoleString(member.Role)
	if roleString == "" {
		return appErr.InternalServerError("internal server error")
	}
	roleID, err := getRoleIDFromDB(roleString)
	if err != nil {
		return err
	}

	if isInserting {
		err = insertChatMemberToDB(tx, member, roleID)
		if err != nil {
			return err
		}
	} else {
		err = updateChatMemberInDB(tx, member, roleID)
		if err != nil {
			return err
		}
	}
	return nil
}
