package authController

import (
	"backend/internal/dtos/userDTO"
	"backend/internal/models/token"
	"backend/internal/models/user"
	"backend/internal/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TODO: придумать как форматы реквестов вынести
// TODO: рефакторинг функции (а то че то она огромная получилась)

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
	_, err = token.InsertToken(userDTO.ID, refreshToken)
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
		"user":        userDTO,
		"accessToken": accessToken,
	})
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

// Activate Refister request format
type ActivateRequest struct {
	UserID uint64 `json:"user_id"`
	Code   string `json:"code"`
}

// Account Activation
func ActivateAccount(c *fiber.Ctx) error {
	var req ActivateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}
	req.trimSpaces()

	err := user.ActivateAccount(req.UserID, req.Code)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "inernal server error" {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "account activated successfully",
	})
}

func (a *ActivateRequest) trimSpaces() {
	a.Code = strings.TrimSpace(a.Code)
}

// Resend activation code request format
type ResendActivationCodeRequest struct {
	UserID uint64 `json:"user_id"`
}

// Resend Activation Code
func ResendActivationCode(c *fiber.Ctx) error {
	var req ResendActivationCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	activationCode, err := user.GetActivationCode(req.UserID)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "inernal server error" {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	activationCode.Regenerate()
	err = activationCode.Save()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	activationCode.SendToEmail()

	return c.JSON(fiber.Map{
		"message": "new code has been sent",
	})
}

// Login request format
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Login
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}
	req.trimSpaces()

	userObject, err := user.GetUserByUsername(req.Login)
	if err != nil && err.Error() == "user not found" {
		userObject, err = user.GetUserByEmail(req.Login)
		if err != nil && err.Error() == "user not found" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid login or password",
			})
		}
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if !userObject.CheckPassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid login or password",
		})
	}

	userDTO := userDTO.CreateUserDTOFromUserObj(userObject)
	accessToken, refreshToken, err := token.GenerateTokens(userDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user has been created, but an error occurred while generating tokens",
		})
	}
	_, err = token.InsertToken(userDTO.ID, refreshToken)
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
		"user":        userDTO,
		"accessToken": accessToken,
	})
}

func (l *LoginRequest) trimSpaces() {
	l.Login = strings.TrimSpace(l.Login)
	l.Password = strings.TrimSpace(l.Password)
}

// Logout
func Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	token.RemoveToken(refreshToken)

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"message": "logout"})
}

// Refresh
func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var userDTO *userDTO.UserDTO
	userDTO, err := token.ValidateRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if userDTO == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	newAccessToken, newRefreshToken, err := token.GenerateTokens(userDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	_, err = token.UpdateToken(refreshToken, newRefreshToken, userDTO.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":        userDTO,
		"accessToken": newAccessToken,
	})
}
