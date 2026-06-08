package socialservice

import (
	filepb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/file/v1"
	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/utils"
)

func toSocialUserProto(a, u *userpb.User, ur *userrelationship.UserRelationship, includePrivacySettings bool, avatarFiles map[string]*filepb.File) *pb.SocialUser {
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
		Avatar:           toSocialUserPhotoProto(u, avatarFiles),
	}
}

func toShortSocialUserProto(a, u *userpb.User, ur *userrelationship.UserRelationship, includePrivacySettings bool, avatarFiles map[string]*filepb.File) *pb.ShortSocialUser {
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
		Avatar:           toSocialUserPhotoProto(u, avatarFiles),
	}
}

// toSocialUserPhotoProto builds a SocialUserPhoto from the user's avatar UserPhoto
// (file_id returned by user-service) and the resolved File from file-service.
// avatarFiles maps file_id -> File; nil map means no resolution available.
func toSocialUserPhotoProto(u *userpb.User, avatarFiles map[string]*filepb.File) *pb.SocialUserPhoto {
	if u.Avatar == nil {
		return nil
	}
	photo := &pb.SocialUserPhoto{
		FileId:    u.Avatar.FileId,
		CreatedAt: u.Avatar.CreatedAt,
	}
	if avatarFiles != nil {
		if f, ok := avatarFiles[u.Avatar.FileId]; ok {
			for _, v := range f.Variants {
				variant := &pb.SocialPhotoVariant{
					VariantType:    v.VariantType,
					UrlForDownload: v.UrlForDownload,
					Size:           v.Size,
					MimeType:       v.MimeType,
				}
				photo.Variants = append(photo.Variants, variant)
			}
		}
	}
	return photo
}

func toSocialUserPhotosProto(userPhotos []*userpb.UserPhoto, fileMap map[string]*filepb.File) []*pb.SocialUserPhoto {
	result := make([]*pb.SocialUserPhoto, 0, len(userPhotos))
	for _, up := range userPhotos {
		photo := &pb.SocialUserPhoto{
			FileId:    up.FileId,
			CreatedAt: up.CreatedAt,
		}
		if fileMap != nil {
			if f, ok := fileMap[up.FileId]; ok {
				for _, v := range f.Variants {
					photo.Variants = append(photo.Variants, &pb.SocialPhotoVariant{
						VariantType:    v.VariantType,
						UrlForDownload: v.UrlForDownload,
						Size:           v.Size,
						MimeType:       v.MimeType,
					})
				}
			}
		}
		result = append(result, photo)
	}
	return result
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
