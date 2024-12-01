package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/chatMember"
	"backend/internal/models/shortUser"
	"backend/internal/models/user"
	"database/sql"
	"errors"
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
		logger.GetInstance().Error(err.Error(), "save chat with members", nil, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	if err := chat.Save(tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	var members []chatMember.ChatMember
	for _, member := range newMembers {
		member.ChatID = chat.ID
		err := member.Save(tx, true)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		members = append(members, member)
	}

	if err := tx.Commit(); err != nil {
		logger.GetInstance().Error(err.Error(), "save chat with members", nil, err)
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
		logger.GetInstance().Error(err.Error(), "rename chat", nil, err)
		return appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	err = chat.Save(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.GetInstance().Error(err.Error(), "rename chat", nil, err)
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
		logger.GetInstance().Error(err.Error(), "delete chat", nil, err)
		return appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	err = chat.Save(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.GetInstance().Error(err.Error(), "delete chat", nil, err)
		return appErr.InternalServerError("internal server error")
	}

	return nil
}

// Get chat members
func (chat *Chat) GetChatMembers(actor *chatMember.ChatMember) ([]chatMember.ChatMember, error) {
	members, err := chatMember.GetChatMembers(actor.User.ID, chat.ID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// Get chat by id from db
func GetChatByID(id uint64) (*Chat, error) {
	chat, err := getChatFromDB(id)
	if err != nil {
		return nil, err
	}

	if chat.Name != nil && *chat.Name != "" {
		chat.IsGroupChat = true
	}

	return chat, nil
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

// save chat in db
func (chat *Chat) Save(tx *sql.Tx) error {
	if chat.ID == 0 {
		err := insertChatToDB(tx, chat)
		if err != nil {
			return err
		}
	} else {
		err := updateChatInDB(tx, chat)
		if err != nil {
			return err
		}
	}
	return nil
}
