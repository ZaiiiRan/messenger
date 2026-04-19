package utils

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user/userDTO"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Parsing date from string
func ParseDate(date string) (*time.Time, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return nil, errors.New("invalid date format (use DD.MM.YYYY)")
	}
	return &parsedDate, nil
}

// Pointer in time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Pointer on string
func StringPtr(s string) *string {
	return &s
}

// Get user dto from locals
func GetUserDTOFromLocals(c *fiber.Ctx) (*userDTO.UserDTO, error) {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return nil, appErr.Unauthorized("unauthorized")
	}
	return user, nil
}
