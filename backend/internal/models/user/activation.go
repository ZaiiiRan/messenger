package user

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/services/mailService"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ActivationCode struct {
	ID        uint64
	UserID    uint64
	Code      string
	ExpiresAt time.Time
}

// Creating activation code object
func CreateActivationCode(userID uint64) *ActivationCode {
	code := &ActivationCode{
		UserID:    userID,
		Code:      strconv.Itoa(generateCode()),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	return code
}

// Get activation code from DataBase
func GetActivationCode(userID uint64) (*ActivationCode, error) {
	db := pgDB.GetDB()
	isActivated, err := IsUserActivated(userID)
	if err != nil {
		return nil, err
	}
	if isActivated {
		return nil, appErr.BadRequest("user already activated")
	}

	var activationCode ActivationCode
	err = db.QueryRow(
		"SELECT id, user_id, code, expires_at FROM activation_codes WHERE user_id = $1", userID,
	).Scan(&activationCode.ID, &activationCode.UserID, &activationCode.Code, &activationCode.ExpiresAt)
	if err == sql.ErrNoRows {
		activationCode = *CreateActivationCode(userID)
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get activation code by userID", userID, err)
		return nil, appErr.InternalServerError("failed to retrieve activation code")
	}
	return &activationCode, nil
}

// Saving activation code in DataBase
func (c *ActivationCode) Save() error {
	db := pgDB.GetDB()
	if c.ID == 0 {
		// new code
		query := `INSERT INTO activation_codes (user_id, code, expires_at) VALUES ($1, $2, $3) RETURNING id`
		err := db.QueryRow(query, c.UserID, c.Code, c.ExpiresAt.UTC()).Scan(&c.ID)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "activation code inserting", c, err)
			return appErr.InternalServerError("internal server error")
		}
	} else {
		// existing code
		query := `UPDATE activation_codes SET code = $1, expires_at = $2 WHERE id = $3`
		_, err := db.Exec(query, c.Code, c.ExpiresAt.UTC(), c.ID)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "activation code updating", c, err)
			return appErr.InternalServerError("internal server error")
		}
	}
	return nil
}

// Deletion activation code from DataBase
func (c *ActivationCode) Delete() error {
	db := pgDB.GetDB()
	if c.ID == 0 {
		return appErr.BadRequest("activation code not found")
	}
	_, err := db.Exec(`DELETE FROM activation_codes WHERE id = $1`, c.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "activation code deleting", c, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// Regenerate code
func (c *ActivationCode) Regenerate() {
	c.Code = strconv.Itoa(generateCode())
	c.ExpiresAt = time.Now().Add(time.Hour)
}

// Sending code to email
func (c *ActivationCode) SendToEmail() error {
	db := pgDB.GetDB()
	var email string
	err := db.QueryRow(`SELECT email FROM users WHERE id = $1 AND is_activated = FALSE`, c.UserID).Scan(&email)
	if err == sql.ErrNoRows {
		return appErr.BadRequest("user not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get email bu user id", c, err)
		return appErr.InternalServerError("internal server error")
	}

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

	err = mailService.SendMail(email, "Account Activation", htmlContent)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "sending activation code to email", c, err)
		return appErr.InternalServerError("there was an error sending the account activation code")
	}
	return nil
}

// Account activation
func ActivateAccount(userID uint64, code string) error {
	activationCode, err := GetActivationCode(userID)
	if err != nil {
		return err
	}
	if activationCode == nil {
		return appErr.BadRequest("activation code not found")
	}

	if time.Now().After(activationCode.ExpiresAt) {
		return appErr.BadRequest("activation code has expired")
	}
	if code != activationCode.Code {
		return appErr.BadRequest("invalid activation code")
	}

	db := pgDB.GetDB()
	_, err = db.Exec(`UPDATE users SET is_activated = TRUE WHERE id = $1`, userID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "account activation", map[string]interface{}{"userID": userID, "code": code}, err)
		return appErr.InternalServerError("failed to activate user account")
	}
	activationCode.Delete()
	return nil
}

// Activation Checking
func IsUserActivated(userID uint64) (bool, error) {
	db := pgDB.GetDB()
	var isActivated bool
	err := db.QueryRow("SELECT is_activated FROM users WHERE id = $1", userID).Scan(&isActivated)
	if err == sql.ErrNoRows {
		return false, appErr.BadRequest("user not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "account activation checking by userID", userID, err)
		return false, appErr.InternalServerError("internal server error")
	}
	return isActivated, nil
}

// generate activation code
func generateCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}
