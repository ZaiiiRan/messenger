package authController

import (
	"backend/internal/dtos/userDTO"
	"backend/internal/models/user"
	"backend/internal/models/token"
	"backend/internal/utils"
	"time"
	"strings"

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

	var birthdate *time.Time
	if req.Birthdate != nil && *req.Birthdate != "" {
		parsedDate, err := utils.ParseDate(*req.Birthdate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		birthdate = parsedDate
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

	userDTO := userDTO.CreateUserDTOFromUserObj(userObject)
	accessToken, refreshToken, err := token.GenerateTokens(userDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user has been created, but an error occurred while generating tokens",
		})
	}
	_, err = token.SaveToken(userDTO.ID, refreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user has been created, but an error occurred while generating tokens",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": userDTO,
		"accessToken": accessToken,
	})
}

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