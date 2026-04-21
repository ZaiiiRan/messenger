package userActivation

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/user"
	"backend/internal/services/mailService"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ActivationCode struct {
	ID        uint64
	User      *user.User
	Code      string
	ExpiresAt time.Time
}

// Creating activation code object
func CreateActivationCode(user *user.User) *ActivationCode {
	code := &ActivationCode{
		User:      user,
		Code:      strconv.Itoa(generateCode()),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	return code
}

// Get activation code from DataBase
func GetActivationCode(user *user.User) (*ActivationCode, error) {
	if user.IsActivated {
		return nil, appErr.BadRequest("user already activated")
	}

	activationCode, err := getActivationCodeFromDB(user.ID)
	if err != nil {
		return nil, err
	}
	if activationCode == nil {
		activationCode = CreateActivationCode(user)
	}
	activationCode.User = user
	return activationCode, nil
}

// Saving activation code in DataBase
func (c *ActivationCode) Save() error {
	if c.ID == 0 {
		// new code
		err := insertActivationCodeToDB(c)
		if err != nil {
			return err
		}
	} else {
		// existing code
		err := updateActivationCodeInDB(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deletion activation code from DataBase
func (c *ActivationCode) Delete() error {
	if c.ID == 0 {
		return appErr.NotFound("activation code not found")
	}

	err := removeActivationCodeFromDB(c)
	if err != nil {
		return err
	}
	return nil
}

// Regenerate code
func (c *ActivationCode) Regenerate() {
	c.Code = strconv.Itoa(generateCode())
	c.ExpiresAt = time.Now().Add(time.Hour)
}

// Account activation
func ActivateAccount(user *user.User, code string) error {
	activationCode, err := GetActivationCode(user)
	if err != nil {
		return err
	}
	if activationCode == nil {
		return appErr.NotFound("activation code not found")
	}

	if time.Now().After(activationCode.ExpiresAt) {
		return appErr.BadRequest("activation code has expired")
	}
	if code != activationCode.Code {
		return appErr.BadRequest("invalid activation code")
	}

	user.IsActivated = true
	err = user.Save()
	if err != nil {
		return err
	}
	activationCode.Delete()
	return nil
}

// Sending code to email
func (c *ActivationCode) SendToEmail() error {
	htmlContent := fmt.Sprintf(
		`<div style="width: 90%%; margin: 0 auto; background-color: rgb(235, 235, 235);
        	border-radius: 25px; padding: 30px 0px; text-align: center;">
        	<h1 style="font-family: sans-serif; font-weight: 900;color:black;">Account Activation</h1>
        	<table style="width: 80%%; margin: auto; background-color: rgb(211, 211, 211); border-radius: 25px; padding: 20px; border-collapse: collapse;">
            	<tr>
                	<td style="padding: 20px; text-align: center;">
                    	<h2 style="font-family: sans-serif; font-weight: 500;color:black;">Enter the code to confirm:</h2>
                    	<h1 style="font-family: sans-serif; font-weight: 900;color:black;">%s</h1>
                	</td>
            	</tr>
        	</table>
    	</div>`, c.Code)

	err := mailService.SendMail(c.User.Email, "Account Activation", htmlContent)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "sending activation code to email", c, err)
		return appErr.InternalServerError("there was an error sending the account activation code")
	}
	return nil
}

// generate activation code
func generateCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}
