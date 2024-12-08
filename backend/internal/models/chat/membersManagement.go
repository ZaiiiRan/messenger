package chat

import (
	"backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/shortUser"
	"backend/internal/models/user"
	"backend/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Get chat member by id from db
func (chat *Chat) GetChatMemberByID(targetID uint64) (*chatMember.ChatMember, error) {
	member, err := chatMember.GetChatMemberByID(targetID, chat.ID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// Leave from chat
func (chat *Chat) LeaveFromChat(member *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}

	if member.IsRemoved() {
		return nil, appErr.Forbidden("you have been removed from the chat")
	}

	member.RemovedBy = &member.User.ID
	member.RemovedAt = utils.TimePtr(time.Now())
	if member.Role != chatMember.Roles.Owner {
		member.Role = chatMember.Roles.Member
	}

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		logger.GetInstance().Error(err.Error(), "leave from chat", nil, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	err = member.Save(tx, false)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.GetInstance().Error(err.Error(), "leave from chat", nil, err)
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
	if member.IsRemoved() {
		return nil, appErr.Forbidden("you have been removed from the chat")
	}

	member.RemovedBy = nil
	member.RemovedAt = nil
	if member.Role != chatMember.Roles.Owner {
		member.Role = chatMember.Roles.Member
	}

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		logger.GetInstance().Error(err.Error(), "return to chat", nil, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	err = member.Save(tx, false)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.GetInstance().Error(err.Error(), "return to chat", nil, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return member, nil
}

// Add members to chat
func (chat *Chat) AddMembers(newMembersIDs []uint64, addingMember *chatMember.ChatMember) ([]chatMember.ChatMember, error) {
	if !chat.IsGroupChat {
		return nil, appErr.BadRequest("chat is not a group chat")
	}
	if err := chat.validateBeforeAddingMembers(len(newMembersIDs)); err != nil {
		return nil, err
	}

	tx, err := pgDB.GetDB().Begin()
	if err != nil {
		logger.GetInstance().Error(err.Error(), "add members to chat", nil, err)
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
		logger.GetInstance().Error(err.Error(), "add members to chat", nil, err)
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

	if target != nil && target.IsLeft() {
		return nil, appErr.BadRequest(fmt.Sprintf("user with id %d left the chat", targetID))
	}

	if target != nil && target.IsRemoved() {
		target, err = chat.addOldMemberToChat(tx, target, addingMember)
		if err != nil {
			return nil, err
		}
		return target, nil
	}

	if target != nil && !target.IsRemoved() {
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
	target.AddedAt = time.Now()
	target.RemovedAt = nil

	err := target.Save(tx, false)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// adding new member to chat
func (chat *Chat) addNewMemberToChat(tx *sql.Tx, user *user.User, addingMember *chatMember.ChatMember) (*chatMember.ChatMember, error) {
	newMember := &chatMember.ChatMember{
		ChatID:  chat.ID,
		User:    shortUser.CreateShortUserFromUser(user),
		AddedBy: addingMember.User.ID,
		Role:    chatMember.Roles.Member,
		AddedAt: time.Now(),
	}
	err := newMember.Save(tx, true)
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
		logger.GetInstance().Error(err.Error(), "remove members from chat", nil, err)
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
		logger.GetInstance().Error(err.Error(), "remove members from chat", nil, err)
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
	if member.IsRemoved() {
		return nil, appErr.BadRequest(fmt.Sprintf("the member with id %d has already been deleted", memberID))
	}

	if member.AddedBy != removingMember.User.ID {
		if removingMember.Role == chatMember.Roles.NotMember {
			logger.GetInstance().Error(fmt.Sprintf("user with id %d in chat with id %d is not member", memberID, member.ChatID), "remove member from chat", nil, err)
			return nil, appErr.InternalServerError("internal server error")
		}

		if removingMember.Role <= member.Role {
			return nil, appErr.Forbidden("trying to delete a member with a higher role")
		}
	} else if member.AddedBy == removingMember.User.ID && member.Role >= removingMember.Role {
		return nil, appErr.BadRequest("trying to delete a member with a higher role")
	}
	member.RemovedBy = &removingMember.User.ID

	if member.IsLeft() {
		member.RemovedAt = utils.TimePtr(time.Now())
	}

	err = member.Save(tx, false)
	if err != nil {
		return nil, err
	}

	return member, nil
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
	if member.IsRemoved() || member.IsLeft() {
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
		logger.GetInstance().Error(err.Error(), "chat member role changing", nil, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	err = member.Save(tx, false)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.GetInstance().Error(err.Error(), "chat member role changing", nil, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return member, nil
}
