package user

import (
	pgDB "backend/internal/dbs/pgDB"
	"backend/internal/services/mailService"
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"time"
	"fmt"
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
		UserID: userID,
		Code: strconv.Itoa(generateCode()),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	return code
}

// Get activation code from DataBase
func GetActivationCode(userID uint64) (*ActivationCode, error) {
	db := pgDB.GetDB()
	var activationCode ActivationCode
	err := db.QueryRow(
		"SELECT id, user_id, code, expires_at FROM activation_codes WHERE user_id = $1", userID,
	).Scan(&activationCode.ID, &activationCode.UserID, &activationCode.Code, &activationCode.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("activation code not found")
	}
	if err != nil {
		return nil, errors.New("failed to retrieve activation code")
	}
	return &activationCode, nil
}

// Saving activation code in DataBase
func (c *ActivationCode) Save() error {
	db := pgDB.GetDB()
	if c.ID == 0 {
		// new code
		query := `INSERT INTO activation_codes (user_id, code, expires_at) VALUES ($1, $2, $3) RETURNING id`
		err := db.QueryRow(query, c.UserID, c.Code, c.ExpiresAt).Scan(&c.ID)
		if err != nil {
			return errors.New("internal server error")
		}
	} else {
		// existing code
		query := `UPDATE activation_codes SET code = $1, expires_at = $2`
		_, err := db.Exec(query, c.Code, c.ExpiresAt)
		if err != nil {
			return errors.New("internal server error")
		}
	}
	return nil
}

// Deletion activation code from DataBase
func (c *ActivationCode) Delete() error {
	db := pgDB.GetDB()
	if c.ID == 0 {
		return errors.New("activation code not found")
	}
	_, err := db.Exec(`DELETE FROM activation_codes WHERE id = $1`, c.ID)
	if err != nil {
		return errors.New("internal server error")
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
	err := db.QueryRow(`SELECT email FROM users WHERE id = $1`, c.UserID).Scan(&email)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	} else if err != nil {
		return errors.New("internal server error")
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
		return errors.New("there was an error sending the account activation code")
	}
	return nil
}

// Account activation
func ActivateAccount(userID uint64, code string) error {
	activationCode, err := GetActivationCode(userID)
	if err != nil {
		return err
	}

	if time.Now().After(activationCode.ExpiresAt) {
		return errors.New("activation code has expired")
	}
	if code != activationCode.Code {
		return errors.New("invalid activation code")
	}

	db := pgDB.GetDB()
	_, err = db.Exec(`UPDATE users SET is_activated = TRUE WHERE id = $1`, userID)
	if err != nil {
		return errors.New("failed to activate user account")
	}
	activationCode.Delete()
	return nil
}

// generate activation code
func generateCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}