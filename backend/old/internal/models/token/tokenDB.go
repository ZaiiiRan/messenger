package token

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"database/sql"
)

// get token by string from db
func getTokenFromDB(refreshToken string) (*Token, error) {
	db := pgDB.GetDB()
	var token Token
	err := db.QueryRow(`SELECT * FROM tokens WHERE refresh_token = $1`, refreshToken).Scan(
		&token.ID, &token.UserID, &token.RefreshToken)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("token not found")
	}
	if err != nil {
		logger.GetInstance().Error(err.Error(), "find token by string", refreshToken, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return &token, nil
}

// insert token to db
func insertTokenToDB(token *Token) error {
	db := pgDB.GetDB()
	err := db.QueryRow(`INSERT INTO tokens (user_id, refresh_token) VALUES ($1, $2) RETURNING *`, token.UserID, token.RefreshToken).Scan(
		&token.ID, &token.UserID, &token.RefreshToken)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token inserting", map[string]interface{}{"userID": token.UserID, "token": token.RefreshToken}, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update token in db
func updateTokenInDB(token *Token) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`UPDATE tokens SET refresh_token = $1 WHERE id = $2 AND user_id = $3 RETURNING *`, token.RefreshToken, token.ID, token.UserID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token updating", map[string]interface{}{"tokenID": token.ID, "userID": token.UserID, "token": token.RefreshToken}, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// remove token from db
func removeTokenFromDB(token *Token) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`DELETE FROM tokens WHERE refresh_token = $1`, token.RefreshToken)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token deleting by string", token.RefreshToken, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// remove other tokens from db
func removeOtherTokensFromDB(token *Token) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`DELETE FROM tokens WHERE user_id = $1 AND id != $2`, token.UserID, token.RefreshToken)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "other tokens deleting by userID", token.UserID, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// remove all tokens from db
func removeAllTokensFromDB(userID uint64) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`DELETE FROM tokens WHERE user_id = $1`, userID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "all tokens deleting by userID", userID, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}
