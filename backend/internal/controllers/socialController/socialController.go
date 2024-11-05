package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/socialUser"
	"backend/internal/models/user"
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

// User manipulate request format
type UserManipulateRequest struct {
	UserID uint64 `json:"user_id"`
}

// Add friend
func AddFriend(c *fiber.Ctx) error {
	var req UserManipulateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}
	err := socialUser.AddFriend(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "ok",
	})
}

// Remove friend
func RemoveFriend(c *fiber.Ctx) error {
	var req UserManipulateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}
	err := socialUser.RemoveFriend(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "friend deleted",
	})
}

// Block user
func BlockUser(c *fiber.Ctx) error {
	var req UserManipulateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}
	err := socialUser.BlockUser(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"status": "user blocked",
	})
}

// Unblock user
func UnblockUser(c *fiber.Ctx) error {
	var req UserManipulateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}
	err := socialUser.UnblockUser(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"status": "user unblocked",
	})
}
