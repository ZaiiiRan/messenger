package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/socialUser"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Fetch users request format
type FetchUsersRequest struct {
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// trim spaces in fetch users request
func (f *FetchUsersRequest) trimSpaces() {
	f.Search = strings.TrimSpace(f.Search)
}

// parsing request body for fetching user requests
func readFetchUsersRequest(c *fiber.Ctx) (*FetchUsersRequest, error) {
	var req FetchUsersRequest
	if err := c.BodyParser(&req); err != nil {
		return &req, appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()
	return &req, nil
}

// fetch user list
func fetchUserList(c *fiber.Ctx, req *FetchUsersRequest, fetchFunc func(userID uint64, search string, limit, offset int) ([]socialUser.SocialUser, error)) error {
	user, err := getUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	users, err := fetchFunc(user.ID, req.Search, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

// Get Users
func GetUsers(c *fiber.Ctx) error {
	req, err := readFetchUsersRequest(c)
	if err != nil {
		return err
	}
	if req.Search == "" {
		return appErr.BadRequest("search parameter is empty")
	} else if len(req.Search) < 4 {
		return appErr.BadRequest("search parameter is very short")
	}

	return fetchUserList(c, req, socialUser.GetUsersByUsernameOrEmail)
}

// Get friends
func GetFriends(c *fiber.Ctx) error {
	req, err := readFetchUsersRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserFriendsByUsernameOrEmail)
}

// Get incoming friend requests
func GetIncomingFriendRequests(c *fiber.Ctx) error {
	req, err := readFetchUsersRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserIncomingFriendRequestsByUsernameOrEmail)
}

// Get outgoing friend requests
func GetOutgoingFriendRequests(c *fiber.Ctx) error {
	req, err := readFetchUsersRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserOutgoingFriendRequestsByUsernameOrEmail)
}

// Get blocked users
func GetBlockedUsers(c *fiber.Ctx) error {
	req, err := readFetchUsersRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserBlockListByUsernameOrEmail)
}
