package chat

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/socialUser"
	"backend/internal/models/user"
	"errors"
	"fmt"
)

// checking user access (friend status and activation) for adding to chat
func checkUserAccess(target *user.User, requestSendingMemberID uint64, isGroup bool) error {
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

	if (relations != nil && *relations == "blocked by target") {
		return appErr.BadRequest(fmt.Sprintf("you are blocked by user with id %d", target.ID))
	}

	if isGroup && ((relations != nil && *relations != "accepted") || (relations == nil)) {
		return appErr.BadRequest(fmt.Sprintf("user with id %d is not in your friends list", target.ID))
	} else if !isGroup && relations != nil && *relations == "blocked"  {
		return appErr.BadRequest(fmt.Sprintf("user with id %d is blocked by you", target.ID))
	}

	return nil
}

// get user object with access checking
func getUserForAdding(userID uint64, requestSendingMemberID uint64, isGroup bool) (*user.User, error) {
	user, err := user.GetUserByID(userID)
	if err != nil {
		var appError *appErr.AppError
		if errors.As(err, &appError) && appError.StatusCode == 404 {
			return nil, appErr.NotFound(fmt.Sprintf("user with id %d not found", userID))
		}
		return nil, err
	}
	err = checkUserAccess(user, requestSendingMemberID, isGroup)
	if err != nil {
		return nil, err
	}
	return user, nil
}