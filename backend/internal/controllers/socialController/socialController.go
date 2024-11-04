package socialController

import (
	"backend/internal/models/socialUser"
	"backend/internal/models/user"
	appErr "backend/internal/errors/appError"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Fetch users request format
type FetchUsersRequest struct {
	Search string `json:"search"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

// trim spaces in fetch users request
func (f *FetchUsersRequest) trimSpaces() {
	f.Search = strings.TrimSpace(f.Search)
}

func readRequest(c *fiber.Ctx) (*FetchUsersRequest, error) {
	var req FetchUsersRequest
	if err := c.BodyParser(&req); err != nil {
		return &req, appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()
	return &req, nil
}

func fetchUserList(c *fiber.Ctx, req *FetchUsersRequest, fetchFunc func(userID uint64, search string, limit, offset int) ([]socialUser.SocialUser, error)) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	users, err := fetchFunc(user.ID, req.Search, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

func GetUsers(c *fiber.Ctx) error {
	req, err := readRequest(c)
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

func GetFriends(c *fiber.Ctx) error {
	req, err := readRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserFriendsByUsernameOrEmail)
}

func GetIncomingFriendRequests(c *fiber.Ctx) error {
	req, err := readRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserIncomingFriendRequestsByUsernameOrEmail)
}

func GetOutgoingFriendRequests(c *fiber.Ctx) error {
	req, err := readRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserOutgoingFriendRequestsByUsernameOrEmail)
}

func GetBlockedUsers(c *fiber.Ctx) error {
	req, err := readRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, socialUser.GetUserBlockListByUsernameOrEmail)
}