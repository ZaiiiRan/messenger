package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
)

func SanitizeGetUserByIdRequest(req *pb.GetUserByIdRequest) {
	if req == nil {
		return
	}
	req.Id = strings.TrimSpace(req.Id)
}

func SanitizeGetUsersByIdsRequest(req *pb.GetUsersByIdsRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeGetUserByUsernameRequest(req *pb.GetUserByUsernameRequest) {
	if req == nil {
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
}

func SanitizeGetUsersByUsernamesRequest(req *pb.GetUsersByUsernamesRequest) {
	if req == nil {
		return
	}
	req.Usernames = sanitizeUniqueArray(sanitizeLowerStringArray(req.Usernames))
}

func SanitizeAddUsersToFriendsRequest(req *pb.AddUsersToFriendsRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeRemoveUsersFromFriendsRequest(req *pb.RemoveUsersFromFriendsRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeBlockUsersRequest(req *pb.BlockUsersRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeUnblockUsersRequest(req *pb.UnblockUsersRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeGetFriendsRequest(req *pb.GetFriendsRequest) {
	if req == nil {
		return
	}
	sanitizeUsersRequest(req.Request)
}

func SanitizeGetIncomingFriendRequestsRequest(req *pb.GetIncomingFriendRequestsRequest) {
	if req == nil {
		return
	}
	sanitizeUsersRequest(req.Request)
}

func SanitizeGetOutgoingFriendRequestsRequest(req *pb.GetOutgoingFriendRequestsRequest) {
	if req == nil {
		return
	}
	sanitizeUsersRequest(req.Request)
}

func SanitizeGetBlockedUsersRequest(req *pb.GetBlockedUsersRequest) {
	if req == nil {
		return
	}
	sanitizeUsersRequest(req.Request)
}

func SanitizeSearchUsersRequest(req *pb.SearchUsersRequest) {
	if req == nil {
		return
	}
	sanitizeUsersRequest(req.Request)
}

func sanitizeUsersRequest(req *pb.UsersRequest) {
	if req == nil {
		return
	}
	req.SearchFilter = sanitizeLowerStringPtr(req.SearchFilter)
	if req.SearchFilter != nil {
		searchFilter := *req.SearchFilter
		searchFilter = strings.Join(strings.Fields(searchFilter), " ")
		req.SearchFilter = &searchFilter
	}
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
	return &trimmed
}

func sanitizeLowerStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.ToLower(strings.TrimSpace(*s))
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func SanitizeUpdateMyPrivacySettingsRequest(req *pb.UpdateMyPrivacySettingsRequest) {
	if req == nil {
		return
	}
	req.Fields = sanitizeUniqueArray(sanitizeStringArray(req.Fields))
	sanitizeUpdateUserPrivacySettings(req.PrivacySettings)
}

func sanitizeUpdateUserPrivacySettings(ps *userpb.UpdateUserPrivacySettings) {
	if ps == nil {
		return
	}
	sanitizeUpdateUserPrivacySetting(ps.Avatar)
	sanitizeUpdateUserPrivacySetting(ps.Photos)
	sanitizeUpdateUserPrivacySetting(ps.PhoneNumber)
	sanitizeUpdateUserPrivacySetting(ps.Email)
	sanitizeUpdateUserPrivacySetting(ps.Birthdate)
	sanitizeUpdateUserPrivacySetting(ps.OnlineStatus)
	sanitizeUpdateUserPrivacySetting(ps.FirstDialogsInit)
	sanitizeUpdateUserPrivacySetting(ps.GroupChatInvites)
}

func sanitizeUpdateUserPrivacySetting(ps *userpb.UpdateUserPrivacySetting) {
	if ps == nil {
		return
	}
	ps.Favourites = sanitizeUniqueArray(sanitizeStringArray(ps.Favourites))
	ps.Exceptions = sanitizeUniqueArray(sanitizeStringArray(ps.Exceptions))
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
