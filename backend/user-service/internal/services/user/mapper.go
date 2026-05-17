package userservice

import (
	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	privacysettings "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/privacy_settings"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
)

func userToProto(u *user.User, includePrivacySettings bool) *pb.User {
	pbUser := &pb.User{
		Id:       u.GetID(),
		Username: u.GetUsername(),
		Email:    u.GetEmail(),
		Profile:  profileToProto(u.GetProfile()),
		Status:   statusToProto(u.GetStatus()),
	}
	if includePrivacySettings {
		pbUser.PrivacySettings = privacySettingsToProto(u.GetPrivacySettings())
	}
	return pbUser
}

func profileToProto(p *profile.Profile) *pb.Profile {
	if p == nil {
		return nil
	}
	return &pb.Profile{
		FirstName: p.GetFirstName(),
		LastName:  p.GetLastName(),
		Phone:     p.GetPhone(),
		Birthdate: utils.FormatDatePtr(p.GetBirthdate()),
		Bio:       p.GetBio(),
	}
}

func statusToProto(s *status.Status) *pb.UserStatus {
	if s == nil {
		return nil
	}
	return &pb.UserStatus{
		IsConfirmed:         s.IsConfirmed(),
		IsPermanentlyBanned: s.IsPermanentlyBanned(),
		BannedUntil:         utils.FormatTimestampPtr(s.GetBannedUntil()),
		IsDeleted:           s.IsDeleted(),
		DeletedAt:           utils.FormatTimestampPtr(s.GetDeletedAt()),
		EmailUpdatedAt:      utils.FormatTimestamp(s.GetEmailUpdatedAt()),
	}
}

func privacySettingsToProto(s *privacysettings.PrivacySettings) *pb.UserPrivacySettings {
	if s == nil {
		return nil
	}
	return &pb.UserPrivacySettings{
		Avatar:           privacySettingToProto(s.GetAvatar()),
		Photos:           privacySettingToProto(s.GetPhotos()),
		PhoneNumber:      privacySettingToProto(s.GetPhoneNumber()),
		Email:            privacySettingToProto(s.GetEmail()),
		Birthdate:        privacySettingToProto(s.GetBirthdate()),
		OnlineStatus:     privacySettingToProto(s.GetOnlineStatus()),
		FirstDialogsInit: privacySettingToProto(s.GetFirstDialogsInit()),
		GroupChatInvites: privacySettingToProto(s.GetGroupChatInvites()),
	}
}

func privacySettingToProto(s privacysettings.PrivacySetting) *pb.UserPrivacySetting {
	return &pb.UserPrivacySetting{
		Value:      s.GetValue().String(),
		Favourites: s.GetFavourites(),
		Exceptions: s.GetExceptions(),
	}
}
