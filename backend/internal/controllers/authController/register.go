package authController

import (
	"backend/internal/models/user"
	"backend/internal/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Register request format
type RegisterRequest struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Phone     *string `json:"phone"`
	Birthdate *string `json:"birthdate"`
}

// Register user
func RegisterUser(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}
	req.trimSpaces()

	birthdate, err := parseBirthdate(req.Birthdate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userObject, err := user.CreateUser(req.Username, req.Email, req.Password, req.Firstname, req.Lastname, req.Phone, birthdate)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "inernal server error" {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := userObject.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	activationCode := user.CreateActivationCode(userObject.ID)
	err = activationCode.Save()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user has been created, but an error occurred while generating activation code",
		})
	}
	activationCode.SendToEmail()

	return createUserDTOAndTokensResponse(userObject, c)
}

// parsing birthdate from request string
func parseBirthdate(date *string) (*time.Time, error) {
	var birthdate *time.Time
	if date != nil && *date != "" {
		parsedDate, err := utils.ParseDate(*date)
		if err != nil {
			return nil, err
		}
		birthdate = parsedDate
	}
	return birthdate, nil
}

// trim spaces in register request
func (r *RegisterRequest) trimSpaces() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.Firstname = strings.TrimSpace(r.Firstname)
	r.Lastname = strings.TrimSpace(r.Lastname)
	if r.Phone != nil {
		trimmedPhone := strings.TrimSpace(*r.Phone)
		r.Phone = &trimmedPhone
	}
	if r.Birthdate != nil {
		trimmedBirthdate := strings.TrimSpace(*r.Birthdate)
		r.Birthdate = &trimmedBirthdate
	}
}