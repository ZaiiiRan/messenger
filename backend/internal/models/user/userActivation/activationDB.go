package userActivation

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"database/sql"
)

// get activation code from db
func getActivationCodeFromDB(userID uint64) (*ActivationCode, error) {
	db := pgDB.GetDB()
	var activationCode ActivationCode
	err := db.QueryRow(
		"SELECT id, code, expires_at FROM activation_codes WHERE user_id = $1", userID,
	).Scan(&activationCode.ID, &activationCode.Code, &activationCode.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get activation code by userID", userID, err)
		return nil, appErr.InternalServerError("failed to retrieve activation code")
	}
	return &activationCode, nil
}

// insert activation code to db
func insertActivationCodeToDB(c *ActivationCode) error {
	db := pgDB.GetDB()
	query := `INSERT INTO activation_codes (user_id, code, expires_at) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, c.User.ID, c.Code, c.ExpiresAt.UTC()).Scan(&c.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "activation code inserting", c, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update activation code in db
func updateActivationCodeInDB(c *ActivationCode) error {
	db := pgDB.GetDB()
	query := `UPDATE activation_codes SET code = $1, expires_at = $2 WHERE id = $3`
	_, err := db.Exec(query, c.Code, c.ExpiresAt.UTC(), c.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "activation code updating", c, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// remove activation code from db
func removeActivationCodeFromDB(c *ActivationCode) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`DELETE FROM activation_codes WHERE id = $1`, c.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "activation code deleting", c, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}
