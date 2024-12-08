package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/shortUser"
	"backend/internal/requests"
	"backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// fetch user list
func fetchUserList(c *fiber.Ctx, req *requests.SearchRequest, fetchFunc func(userID uint64, search string, limit, offset int) ([]shortUser.ShortUser, error)) error {
	user, err := utils.GetUserDTOFromLocals(c)
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

// Get Users (id, username, firstname and lastname)
func GetUsers(c *fiber.Ctx) error {
	req, err := requests.ParseSearchRequest(c)
	if err != nil {
		return err
	}
	if req.Search == "" {
		return appErr.BadRequest("search parameter is empty")
	} else if len(req.Search) < 4 {
		return appErr.BadRequest("search parameter is very short")
	}

	return fetchUserList(c, req, shortUser.SearchAll)
}

// Get friends
func GetFriends(c *fiber.Ctx) error {
	req, err := requests.ParseSearchRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, shortUser.SearchFriends)
}

// Get incoming friend requests
func GetIncomingFriendRequests(c *fiber.Ctx) error {
	req, err := requests.ParseSearchRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, shortUser.SearchIncomingFriendRequests)
}

// Get outgoing friend requests
func GetOutgoingFriendRequests(c *fiber.Ctx) error {
	req, err := requests.ParseSearchRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, shortUser.SearchOutgoingFriendRequests)
}

// Get blocked users
func GetBlockedUsers(c *fiber.Ctx) error {
	req, err := requests.ParseSearchRequest(c)
	if err != nil {
		return err
	}
	return fetchUserList(c, req, shortUser.SearchBlockList)
}
