package userrelationship

import "time"

type UserRelationship struct {
	userId1   string
	userId2   string
	status    int16
	createdAt time.Time
	updatedAt time.Time
}

func New(userId1, userId2 string, status UserRelationshipStatus) *UserRelationship {
	now := time.Now()

	if userId1 > userId2 {
		userId1, userId2 = userId2, userId1
	}

	return &UserRelationship{
		userId1:   userId1,
		userId2:   userId2,
		status:    int16(status),
		createdAt: now,
		updatedAt: now,
	}
}

func FromStorage(userId1, userId2 string, status int16, createdAt, updatedAt time.Time) *UserRelationship {
	return &UserRelationship{
		userId1:   userId1,
		userId2:   userId2,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (ur *UserRelationship) GetUserID1() string      { return ur.userId1 }
func (ur *UserRelationship) GetUserID2() string      { return ur.userId2 }
func (ur *UserRelationship) GetCreatedAt() time.Time { return ur.createdAt }
func (ur *UserRelationship) GetUpdatedAt() time.Time { return ur.updatedAt }

func (ur *UserRelationship) GetStatus() UserRelationshipStatus {
	return UserRelationshipStatus(ur.status)
}

func (ur *UserRelationship) RoleOf(userID string) int {
	if ur.userId1 == userID {
		return 1
	}
	return 2
}

func (ur *UserRelationship) SetStatus(status UserRelationshipStatus) error {
	current := UserRelationshipStatus(ur.status)

	if status == Friends && current != FriendRequestBy1 && current != FriendRequestBy2 {
		return ErrCannotBecomeFriends
	}
	if status == BlockedByBoth && current != BlockedBy1 && current != BlockedBy2 {
		return ErrCannotBeMutualBlock
	}

	ur.status = int16(status)
	return nil
}

func (ur *UserRelationship) SetUpdatedAt(now *time.Time) {
	if now == nil {
		t := time.Now()
		now = &t
	}
	ur.updatedAt = *now
}
