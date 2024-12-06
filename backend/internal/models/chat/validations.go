package chat

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user/userDTO"
	"errors"
	"fmt"
)

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
func validateBeforeCreatingChat(name string, members []uint64, isGroup bool, ownerDTO *userDTO.UserDTO) error {
	if isGroup {
		if err := validateName(name); err != nil {
			return err
		}
		if len(members) < 2 {
			return appErr.BadRequest("need at least 2 members for group chat")
		}
		if len(members) > 1000 {
			return appErr.BadRequest("maximum number of chat members: 1000")
		}
	} else {
		if len(members) < 1 {
			return appErr.BadRequest("need at least 1 member for private chat")
		} else if len(members) > 1 {
			return appErr.BadRequest("max 1 member for private chat")
		}

		privateChatExists, err := checkPrivateChatExists(ownerDTO.ID, members[0])
		if err != nil {
			return err
		}
		if privateChatExists {
			return appErr.BadRequest(fmt.Sprintf("chat between user with id %d and %d already exists", ownerDTO.ID, members[0]))
		}
	}

	return nil
}

// validate chat members count before adding
func (chat *Chat) validateBeforeAddingMembers(addingCount int) error {
	count, err := chat.GetChatMembersCount()
	if err != nil {
		return err
	}
	if (count + addingCount) > 1000 {
		return appErr.BadRequest("maximum number of chat members: 1000")
	}
	return nil
}

// chat existing checking
func checkPrivateChatExists(member1, member2 uint64) (bool, error) {
	var appError *appErr.AppError
	chat, err := getPrivateChatFromDB(member1, member2)
	if err != nil && errors.As(err, &appError) && appError.StatusCode == 404 {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return chat.ID == 0, nil
}
