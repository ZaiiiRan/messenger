package handlers

import (
	socialpb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/social/v1"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/httputil"
	"github.com/gofiber/fiber/v2"
)

// usersRequestFromQuery builds the shared paging/search request used by the
// social listing endpoints from the query string.
func usersRequestFromQuery(c *fiber.Ctx) *socialpb.UsersRequest {
	req := &socialpb.UsersRequest{
		Page:                   int32(c.QueryInt("page", 0)),
		PageSize:               int32(c.QueryInt("page_size", 0)),
		IncludePrivacySettings: c.QueryBool("include_privacy_settings", false),
		SortByUpdatedAt:        c.QueryBool("sort_by_updated_at", false),
	}
	if filter := c.Query("search_filter"); filter != "" {
		req.SearchFilter = &filter
	}
	return req
}

// GetSocialUserById godoc
// @Summary      Get a social user by id
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User id"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Success      200 {object} socialpb.GetUserByIdResponse
// @Router       /api/social/users/{id} [get]
func (h *Handlers) GetSocialUserById(c *fiber.Ctx) error {
	req := &socialpb.GetUserByIdRequest{
		Id:                     c.Params("id"),
		IncludePrivacySettings: c.QueryBool("include_privacy_settings", false),
	}

	resp, err := h.clients.Social.GetUserById(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetSocialUsersByIds godoc
// @Summary      Get social users by ids
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.GetUsersByIdsRequest true "User ids"
// @Success      200 {object} socialpb.GetUsersByIdsResponse
// @Router       /api/social/users/by-ids [post]
func (h *Handlers) GetSocialUsersByIds(c *fiber.Ctx) error {
	req := &socialpb.GetUsersByIdsRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.GetUsersByIds(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetSocialUserByUsername godoc
// @Summary      Get a social user by username
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Success      200 {object} socialpb.GetUserByUsernameResponse
// @Router       /api/social/users/username/{username} [get]
func (h *Handlers) GetSocialUserByUsername(c *fiber.Ctx) error {
	req := &socialpb.GetUserByUsernameRequest{
		Username:               c.Params("username"),
		IncludePrivacySettings: c.QueryBool("include_privacy_settings", false),
	}

	resp, err := h.clients.Social.GetUserByUsername(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetSocialUsersByUsernames godoc
// @Summary      Get social users by usernames
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.GetUsersByUsernamesRequest true "Usernames"
// @Success      200 {object} socialpb.GetUsersByUsernamesResponse
// @Router       /api/social/users/by-usernames [post]
func (h *Handlers) GetSocialUsersByUsernames(c *fiber.Ctx) error {
	req := &socialpb.GetUsersByUsernamesRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.GetUsersByUsernames(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// AddUsersToFriends godoc
// @Summary      Add users to friends
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.AddUsersToFriendsRequest true "User ids"
// @Success      200 {object} socialpb.AddUsersToFriendsResponse
// @Router       /api/social/friends [post]
func (h *Handlers) AddUsersToFriends(c *fiber.Ctx) error {
	req := &socialpb.AddUsersToFriendsRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.AddUsersToFriends(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// RemoveUsersFromFriends godoc
// @Summary      Remove users from friends
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.RemoveUsersFromFriendsRequest true "User ids"
// @Success      200 {object} socialpb.RemoveUsersFromFriendsResponse
// @Router       /api/social/friends [delete]
func (h *Handlers) RemoveUsersFromFriends(c *fiber.Ctx) error {
	req := &socialpb.RemoveUsersFromFriendsRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.RemoveUsersFromFriends(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// BlockUsers godoc
// @Summary      Block users
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.BlockUsersRequest true "User ids"
// @Success      200 {object} socialpb.BlockUsersResponse
// @Router       /api/social/blocked [post]
func (h *Handlers) BlockUsers(c *fiber.Ctx) error {
	req := &socialpb.BlockUsersRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.BlockUsers(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// UnblockUsers godoc
// @Summary      Unblock users
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.UnblockUsersRequest true "User ids"
// @Success      200 {object} socialpb.UnblockUsersResponse
// @Router       /api/social/blocked [delete]
func (h *Handlers) UnblockUsers(c *fiber.Ctx) error {
	req := &socialpb.UnblockUsersRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.UnblockUsers(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetFriends godoc
// @Summary      List friends
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        search_filter query string false "Search filter"
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Param        sort_by_updated_at query bool false "Sort by updated at"
// @Success      200 {object} socialpb.GetFriendsResponse
// @Router       /api/social/friends [get]
func (h *Handlers) GetFriends(c *fiber.Ctx) error {
	req := &socialpb.GetFriendsRequest{Request: usersRequestFromQuery(c)}

	resp, err := h.clients.Social.GetFriends(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetIncomingFriendRequests godoc
// @Summary      List incoming friend requests
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        search_filter query string false "Search filter"
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Param        sort_by_updated_at query bool false "Sort by updated at"
// @Success      200 {object} socialpb.GetIncomingFriendRequestsResponse
// @Router       /api/social/friends/incoming [get]
func (h *Handlers) GetIncomingFriendRequests(c *fiber.Ctx) error {
	req := &socialpb.GetIncomingFriendRequestsRequest{Request: usersRequestFromQuery(c)}

	resp, err := h.clients.Social.GetIncomingFriendRequests(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetOutgoingFriendRequests godoc
// @Summary      List outgoing friend requests
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        search_filter query string false "Search filter"
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Param        sort_by_updated_at query bool false "Sort by updated at"
// @Success      200 {object} socialpb.GetOutgoingFriendRequestsResponse
// @Router       /api/social/friends/outgoing [get]
func (h *Handlers) GetOutgoingFriendRequests(c *fiber.Ctx) error {
	req := &socialpb.GetOutgoingFriendRequestsRequest{Request: usersRequestFromQuery(c)}

	resp, err := h.clients.Social.GetOutgoingFriendRequests(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetBlockedUsers godoc
// @Summary      List blocked users
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        search_filter query string false "Search filter"
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Param        sort_by_updated_at query bool false "Sort by updated at"
// @Success      200 {object} socialpb.GetBlockedUsersResponse
// @Router       /api/social/blocked [get]
func (h *Handlers) GetBlockedUsers(c *fiber.Ctx) error {
	req := &socialpb.GetBlockedUsersRequest{Request: usersRequestFromQuery(c)}

	resp, err := h.clients.Social.GetBlockedUsers(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// SearchUsers godoc
// @Summary      Search users
// @Tags         social
// @Produce      json
// @Security     BearerAuth
// @Param        search_filter query string false "Search filter"
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Param        sort_by_updated_at query bool false "Sort by updated at"
// @Success      200 {object} socialpb.SearchUsersResponse
// @Router       /api/social/search [get]
func (h *Handlers) SearchUsers(c *fiber.Ctx) error {
	req := &socialpb.SearchUsersRequest{Request: usersRequestFromQuery(c)}

	resp, err := h.clients.Social.SearchUsers(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// UpdateSocialPrivacySettings godoc
// @Summary      Update my privacy settings (social)
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body socialpb.UpdateMyPrivacySettingsRequest true "Privacy settings"
// @Success      200 {object} socialpb.UpdateMyPrivacySettingsResponse
// @Router       /api/social/me/privacy [patch]
func (h *Handlers) UpdateSocialPrivacySettings(c *fiber.Ctx) error {
	req := &socialpb.UpdateMyPrivacySettingsRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Social.UpdateMyPrivacySettings(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}
