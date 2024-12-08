package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"backend/internal/models/user/userActivation"
	"backend/internal/requests"
	"backend/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Register user
func RegisterUser(c *fiber.Ctx) error {
	req, err := requests.ParseRegisterRequest(c)
	if err != nil {
		return err
	}

	birthdate, err := parseBirthdate(req.Birthdate)
	if err != nil {
		return appErr.BadRequest(err.Error())
	}

	userObject, err := user.CreateUser(req.Username, req.Email, req.Password, req.Firstname, req.Lastname, req.Phone, birthdate)
	if err != nil {
		return err
	}

	if err := userObject.Save(); err != nil {
		return err
	}

	activationCode := userActivation.CreateActivationCode(userObject)
	err = activationCode.Save()
	if err != nil {
		return appErr.InternalServerError("error occured while sending activation code")
	}
	go activationCode.SendToEmail()

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
