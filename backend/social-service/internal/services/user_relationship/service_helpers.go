package userrelationshipservice

import (
	"time"

	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
)

func buildRelationshipMap(urs []*userrelationship.UserRelationship, actorID string) map[string]*userrelationship.UserRelationship {
	m := make(map[string]*userrelationship.UserRelationship, len(urs))
	for _, ur := range urs {
		m[ur.OtherUserID(actorID)] = ur
	}
	return m
}

func applyAddFriend(
	actor *userpb.User,
	friendCandidate *userpb.User,
	existing *userrelationship.UserRelationship,
	now time.Time,
) (*userrelationship.UserRelationship, bool, error) {
	if existing == nil {
		var status userrelationship.UserRelationshipStatus
		if actor.Id > friendCandidate.Id {
			status = userrelationship.FriendRequestBy2
		} else {
			status = userrelationship.FriendRequestBy1
		}
		return userrelationship.New(actor.Id, friendCandidate.Id, status), false, nil
	}

	actorRole := existing.RoleOf(actor.Id)
	curStatus := existing.GetStatus()

	var newStatus userrelationship.UserRelationshipStatus
	blocked := false

	switch {
	case curStatus == userrelationship.Friends:
		return existing, false, ErrAlreadyFriends

	case (curStatus == userrelationship.BlockedBy1 && actorRole == 2) ||
		(curStatus == userrelationship.BlockedBy2 && actorRole == 1):
		return existing, false, ErrBlockedByFriendCandidate

	case (curStatus == userrelationship.FriendRequestBy1 && actorRole == 2) ||
		(curStatus == userrelationship.FriendRequestBy2 && actorRole == 1):
		newStatus = userrelationship.Friends

	case curStatus == userrelationship.FriendRequestBy1 || curStatus == userrelationship.FriendRequestBy2:
		return existing, false, ErrFriendRequestAlreadySent

	case (curStatus == userrelationship.BlockedBy1 && actorRole == 1) ||
		(curStatus == userrelationship.BlockedBy2 && actorRole == 2):
		if actorRole == 1 {
			newStatus = userrelationship.FriendRequestBy1
		} else {
			newStatus = userrelationship.FriendRequestBy2
		}

	case curStatus == userrelationship.BlockedByBoth:
		blocked = true
		if actorRole == 1 {
			newStatus = userrelationship.BlockedBy2
		} else {
			newStatus = userrelationship.BlockedBy1
		}
	}

	if err := existing.SetStatus(newStatus); err != nil {
		return existing, false, ErrAddUserToFriends
	}
	existing.SetUpdatedAt(&now)
	return existing, blocked, nil
}

func applyRemoveFriend(existing *userrelationship.UserRelationship) (*userrelationship.UserRelationship, bool) {
	if existing == nil {
		return nil, false
	}
	s := existing.GetStatus()
	if s == userrelationship.BlockedBy1 || s == userrelationship.BlockedBy2 || s == userrelationship.BlockedByBoth {
		return existing, false
	}
	return existing, true
}

func applyBlockUser(
	actor *userpb.User,
	blockCandidate *userpb.User,
	existing *userrelationship.UserRelationship,
	now time.Time,
) (*userrelationship.UserRelationship, error) {
	if existing == nil {
		var status userrelationship.UserRelationshipStatus
		if actor.Id > blockCandidate.Id {
			status = userrelationship.BlockedBy2
		} else {
			status = userrelationship.BlockedBy1
		}
		return userrelationship.New(actor.Id, blockCandidate.Id, status), nil
	}

	actorRole := existing.RoleOf(actor.Id)
	curStatus := existing.GetStatus()

	var newStatus userrelationship.UserRelationshipStatus
	switch {
	case (curStatus == userrelationship.BlockedBy1 && actorRole == 1) ||
		(curStatus == userrelationship.BlockedBy2 && actorRole == 2) ||
		curStatus == userrelationship.BlockedByBoth:
		return existing, ErrAlreadyBlocked

	case (curStatus == userrelationship.BlockedBy1 && actorRole == 2) ||
		(curStatus == userrelationship.BlockedBy2 && actorRole == 1):
		newStatus = userrelationship.BlockedByBoth

	default:
		if actorRole == 1 {
			newStatus = userrelationship.BlockedBy1
		} else {
			newStatus = userrelationship.BlockedBy2
		}
	}

	if err := existing.SetStatus(newStatus); err != nil {
		return existing, ErrBlockUser
	}
	existing.SetUpdatedAt(&now)
	return existing, nil
}

func applyUnblockUser(
	actor *userpb.User,
	existing *userrelationship.UserRelationship,
	now time.Time,
) (*userrelationship.UserRelationship, bool, bool) {
	if existing == nil {
		return nil, false, true
	}

	curStatus := existing.GetStatus()
	actorRole := existing.RoleOf(actor.Id)

	switch {
	case (curStatus == userrelationship.BlockedBy1 && actorRole == 2) ||
		(curStatus == userrelationship.BlockedBy2 && actorRole == 1) ||
		curStatus == userrelationship.FriendRequestBy1 || curStatus == userrelationship.FriendRequestBy2 ||
		curStatus == userrelationship.Friends:
		return existing, false, true

	case (curStatus == userrelationship.BlockedBy1 && actorRole == 1) ||
		(curStatus == userrelationship.BlockedBy2 && actorRole == 2):
		return existing, true, false

	default:
		var newStatus userrelationship.UserRelationshipStatus
		if actorRole == 1 {
			newStatus = userrelationship.BlockedBy2
		} else {
			newStatus = userrelationship.BlockedBy1
		}
		if err := existing.SetStatus(newStatus); err != nil {
			return existing, false, true
		}
		existing.SetUpdatedAt(&now)
		return existing, false, false
	}
}
