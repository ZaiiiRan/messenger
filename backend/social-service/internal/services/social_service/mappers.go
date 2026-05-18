package socialservice

import (
	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/utils"
)

func toSocialUserProto(a, u *userpb.User, ur *userrelationship.UserRelationship, includePrivacySettings bool) *pb.SocialUser {
	if a == nil || u == nil {
		return nil
	}

	var ps *userpb.UserPrivacySettings
	if includePrivacySettings {
		ps = u.PrivacySettings
	}

	return &pb.SocialUser{
		Id:               u.Id,
		Username:         u.Username,
		Profile:          toSocialUserProfileProto(u),
		Status:           toSocialUserStatusProto(u),
		UserRelationship: toSocialUserRelationshipProto(a, ur),
		PrivacySettings:  ps,
	}
}

func toShortSocialUserProto(a, u *userpb.User, ur *userrelationship.UserRelationship, includePrivacySettings bool) *pb.ShortSocialUser {
	if a == nil || u == nil {
		return nil
	}

	var ps *userpb.UserPrivacySettings
	if includePrivacySettings {
		ps = u.PrivacySettings
	}

	return &pb.ShortSocialUser{
		Id:               u.Id,
		Username:         u.Username,
		Profile:          toShortSocialUserProfileProto(u),
		Status:           toSocialUserStatusProto(u),
		UserRelationship: toSocialUserRelationshipProto(a, ur),
		PrivacySettings:  ps,
	}
}

func toSocialUserRelationshipProto(actor *userpb.User, ur *userrelationship.UserRelationship) string {
	if ur == nil {
		return "no_relation"
	}

	actorRole := ur.RoleOf(actor.Id)
	status := ur.GetStatus()

	switch {
	case status == userrelationship.Friends:
		return "friends"
	case status == userrelationship.FriendRequestBy1 && actorRole == 1 || status == userrelationship.FriendRequestBy2 && actorRole == 2:
		return "outgoing_friend_request"
	case status == userrelationship.FriendRequestBy1 && actorRole == 2 || status == userrelationship.FriendRequestBy2 && actorRole == 1:
		return "incoming_friend_request"
	case status == userrelationship.BlockedBy1 && actorRole == 1 || status == userrelationship.BlockedBy2 && actorRole == 2:
		return "blocked"
	case status == userrelationship.BlockedBy1 && actorRole == 2 || status == userrelationship.BlockedBy2 && actorRole == 1:
		return "blocked_by_target"
	case status == userrelationship.BlockedByBoth:
		return "blocked_by_both"
	}

	return "no_relation"
}

func toSocialUserProfileProto(u *userpb.User) *pb.SocialUserProfile {
	p := &pb.SocialUserProfile{
		Email: utils.StringPtr(u.Email),
	}
	if u.Profile != nil {
		p.FirstName = u.Profile.FirstName
		p.LastName = u.Profile.LastName
		p.Phone = u.Profile.Phone
		p.Birthdate = u.Profile.Birthdate
		p.Bio = u.Profile.Bio
	}
	return p
}

func toShortSocialUserProfileProto(u *userpb.User) *pb.ShortSocialUserProfile {
	if u.Profile == nil {
		return &pb.ShortSocialUserProfile{}
	}
	return &pb.ShortSocialUserProfile{
		FirstName: u.Profile.FirstName,
		LastName:  u.Profile.LastName,
	}
}

func toSocialUserStatusProto(u *userpb.User) *pb.SocialUserStatus {
	if u.Status == nil {
		return &pb.SocialUserStatus{}
	}
	return &pb.SocialUserStatus{
		IsPermanentlyBanned: u.Status.IsPermanentlyBanned,
		IsDeleted:           u.Status.IsDeleted,
	}
}
