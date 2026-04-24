package userservice

import (
	"time"

	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
)

func userToProto(u *user.User) *pb.User {
	return &pb.User{
		Id:       u.GetID(),
		Username: u.GetUsername(),
		Email:    u.GetEmail(),
		Profile:  profileToProto(u.GetProfile()),
		Status:   statusToProto(u.GetStatus()),
	}
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
	var bannedUntil *string
	if bu := s.GetBannedUntil(); bu != nil {
		formatted := bu.Format(time.RFC3339)
		bannedUntil = &formatted
	}
	var deletedAt *string
	if da := s.GetDeletedAt(); da != nil {
		formatted := da.Format(time.RFC3339)
		deletedAt = &formatted
	}
	return &pb.UserStatus{
		IsConfirmed:         s.IsConfirmed(),
		IsPermanentlyBanned: s.IsPermanentlyBanned(),
		BannedUntil:         bannedUntil,
		IsDeleted:           s.IsDeleted(),
		DeletedAt:           deletedAt,
	}
}
