package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
)

func SanitizeCreateUserRequest(req *pb.CreateUserRequest) {
	if req == nil {
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	sanitizeProfile(req.Profile)
}

func SanitizeGetUsersRequest(req *pb.GetUsersRequest) {
	if req == nil {
		return
	}

	req.Ids = sanitizeStringArray(sanitizeUniqueArray(req.Ids))
	req.ExcludeIds = sanitizeStringArray(sanitizeUniqueArray(req.ExcludeIds))
	req.FullUsernames = sanitizeLowerStringArray(sanitizeUniqueArray(req.FullUsernames))
	req.PartialUsernames = sanitizeLowerStringArray(sanitizeUniqueArray(req.PartialUsernames))
	req.FullEmails = sanitizeLowerStringArray(sanitizeUniqueArray(req.FullEmails))
	req.PartialEmails = sanitizeLowerStringArray(sanitizeUniqueArray(req.PartialEmails))
	req.PhoneNumbers = sanitizeStringArray(sanitizeUniqueArray(req.PhoneNumbers))
	req.PartialNames = sanitizeStringArray(sanitizeUniqueArray(req.PartialNames))

	req.SearchFilter = sanitizeStringPtr(req.SearchFilter)
	if req.SearchFilter != nil {
		searchFilter := *req.SearchFilter
		searchFilter = strings.Join(strings.Fields(searchFilter), " ")
		req.SearchFilter = &searchFilter
	}

	req.DeletedFrom = sanitizeStringPtr(req.DeletedFrom)
	req.DeletedTo = sanitizeStringPtr(req.DeletedTo)
	req.CreatedFrom = sanitizeStringPtr(req.CreatedFrom)
	req.CreatedTo = sanitizeStringPtr(req.CreatedTo)
	req.UpdatedFrom = sanitizeStringPtr(req.UpdatedFrom)
	req.UpdatedTo = sanitizeStringPtr(req.UpdatedTo)
}

func SanitizeConfirmUserRequest(req *pb.ConfirmUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeGetUserByIDRequest(req *pb.GetUserByIDRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeGetUserByUsernameRequest(req *pb.GetUserByUsernameRequest) {
	if req == nil {
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
}

func SanitizeGetUserByEmailRequest(req *pb.GetUserByEmailRequest) {
	if req == nil {
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
}

func SanitizeBanUserRequest(req *pb.BanUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
	req.BannedUntil = sanitizeStringPtr(req.BannedUntil)
}

func SanitizeUnbanUserRequest(req *pb.UnbanUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeUpdateUserEmailRequest(req *pb.UpdateUserEmailRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
	req.NewEmail = strings.ToLower(strings.TrimSpace(req.NewEmail))
}

func SanitizeUpdateMeByUserRequest(req *pb.UpdateMeByUserRequest) {
	if req == nil {
		return
	}
	if req.Fields == nil {
		return
	}
	req.Fields = sanitizeUniqueArray(req.Fields)
	sanitizeUpdateUser(req.User)
}

func SanitizeUpdateMyPrivacySettingsByUserRequest(req *pb.UpdateMyPrivacySettingsByUserRequest) {
	if req == nil {
		return
	}
	if req.Fields == nil {
		return
	}
	req.Fields = sanitizeUniqueArray(req.Fields)
	sanitizeUpdateUserPrivacySettings(req.PrivacySettings)
}

func sanitizeProfile(p *pb.Profile) {
	if p == nil {
		return
	}
	p.FirstName = strings.TrimSpace(p.FirstName)
	p.LastName = strings.TrimSpace(p.LastName)
	p.Phone = sanitizeStringPtr(p.Phone)
	p.Birthdate = sanitizeStringPtr(p.Birthdate)
	p.Bio = sanitizeStringPtr(p.Bio)
}

func sanitizeUpdateUser(u *pb.UpdateUser) {
	if u == nil {
		return
	}
	u.Username = sanitizeLowerStringPtr(u.Username)
	sanitizeUpdateProfile(u.Profile)
}

func sanitizeUpdateProfile(p *pb.UpdateProfile) {
	if p == nil {
		return
	}
	p.FirstName = sanitizeStringPtr(p.FirstName)
	p.LastName = sanitizeStringPtr(p.LastName)
	p.Phone = sanitizeStringPtr(p.Phone)
	p.Birthdate = sanitizeStringPtr(p.Birthdate)
	p.Bio = sanitizeStringPtr(p.Bio)
}

func sanitizeUpdateUserPrivacySettings(s *pb.UpdateUserPrivacySettings) {
	if s == nil {
		return
	}
	sanitizeUpdateUserPrivacySetting(s.Avatar)
	sanitizeUpdateUserPrivacySetting(s.Photos)
	sanitizeUpdateUserPrivacySetting(s.PhoneNumber)
	sanitizeUpdateUserPrivacySetting(s.Email)
	sanitizeUpdateUserPrivacySetting(s.Birthdate)
	sanitizeUpdateUserPrivacySetting(s.OnlineStatus)
	sanitizeUpdateUserPrivacySetting(s.FirstDialogsInit)
	sanitizeUpdateUserPrivacySetting(s.GroupChatInvites)
}

func sanitizeUpdateUserPrivacySetting(s *pb.UpdateUserPrivacySetting) {
	if s == nil {
		return
	}
	s.Value = sanitizeStringPtr(s.Value)
	s.Favourites = sanitizeStringArray(sanitizeUniqueArray(s.Favourites))
	s.Exceptions = sanitizeStringArray(sanitizeUniqueArray(s.Exceptions))
}

func sanitizeStringArray(arr []string) []string {
	var sanitized []string
	for _, str := range arr {
		str = strings.TrimSpace(str)
		if str != "" {
			sanitized = append(sanitized, str)
		}
	}
	return sanitized
}

func sanitizeLowerStringArray(arr []string) []string {
	var sanitized []string
	for _, str := range arr {
		str = strings.ToLower(strings.TrimSpace(str))
		if str != "" {
			sanitized = append(sanitized, str)
		}
	}
	return sanitized
}

func sanitizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func sanitizeLowerStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.ToLower(strings.TrimSpace(*s))
	return &trimmed
}

func sanitizeUniqueArray[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	out := s[:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
