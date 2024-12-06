package chatMember

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/shortUser"
	"backend/internal/models/user"
	"database/sql"
	"fmt"
	"time"
)

type ChatMember struct {
	User      *shortUser.ShortUser `json:"user"`
	Role      int                  `json:"role"`
	ChatID    uint64               `json:"chat_id"`
	RemovedBy *uint64
	AddedBy   uint64
	AddedAt   time.Time
}

// Removed checking
func (member *ChatMember) IsRemoved() bool {
	return member.RemovedBy != nil && *member.RemovedBy != member.User.ID
}

func (member *ChatMember) IsLeft() bool {
	return member.RemovedBy != nil && *member.RemovedBy == member.User.ID
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

	removedBy, addedBy, addedAt, err := getChatMemberRemoveAndAddInfo(targetID, chatID)
	if err != nil {
		return nil, err
	}

	member.AddedBy = addedBy
	member.RemovedBy = removedBy
	member.AddedAt = *addedAt

	return &member, nil
}

// Save chat member
func (member *ChatMember) Save(tx *sql.Tx, isInserting bool) error {
	roleString := GetRoleString(member.Role)
	if roleString == "" {
		logger.GetInstance().Error("role string is empty", "save chat member", nil, nil)
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

// Get chat members by search string
func GetChatMembers(actorID, chatID uint64) ([]ChatMember, error) {
	return getChatMembersFromDB(actorID, chatID)
}

// Get chat members count in chat
func GetChatMembersCount(chatID uint64) (int, error) {
	return getChatMembersCountFromDB(chatID)
}
