package token

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/user"
	"backend/internal/utils"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessKey  = os.Getenv("JWT_ACCESS_KEY")
	refreshKey = os.Getenv("JWT_REFRESH_KEY")
)

type Token struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// Find token in DataBase
func FindToken(refreshToken string) (*Token, error) {
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

// Insert token in DataBase
func InsertToken(userID uint64, refreshToken string) (*Token, error) {
	db := pgDB.GetDB()
	var token Token
	err := db.QueryRow(`INSERT INTO tokens (user_id, refresh_token) VALUES ($1, $2) RETURNING *`, userID, refreshToken).Scan(
		&token.ID, &token.UserID, &token.RefreshToken)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token inserting", map[string]interface{}{"userID": userID, "token": refreshToken}, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return &token, nil
}

// Update token in DataBase
func UpdateToken(oldRefreshToken, newRefreshToken string, userID uint64) (*Token, error) {
	db := pgDB.GetDB()
	var token Token
	err := db.QueryRow(`UPDATE tokens SET refresh_token = $1 WHERE refresh_token = $2 AND user_id = $3 RETURNING *`, newRefreshToken, oldRefreshToken, userID).Scan(
		&token.ID, &token.UserID, &token.RefreshToken)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token updating", map[string]interface{}{"userID": userID, "newToken": newRefreshToken, "oldToken": oldRefreshToken}, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return &token, nil
}

// Remove token from DataBase
func RemoveToken(refreshToken string) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`DELETE FROM tokens WHERE refresh_token = $1`, refreshToken)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token deleting by string", refreshToken, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// Generate token pair
func GenerateTokens(payload *user.UserDTO) (string, string, error) {
	// 30 minutes
	accessToken, err := createToken(payload, 30, accessKey)
	if err != nil {
		return "", "", err
	}

	// 30 days
	refreshToken, err := createToken(payload, 43200, refreshKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Validating access token
func ValidateAccessToken(tokenString string) (*user.UserDTO, error) {
	userDTO, _, err := validateToken(tokenString, accessKey)
	return userDTO, err
}

// Validating refresh token
func ValidateRefreshToken(tokenString string) (*user.UserDTO, error) {
	_, err := FindToken(tokenString)
	if err != nil && err.Error() == "token not found" {
		return nil, appErr.Unauthorized("unauthorized")
	} else if err != nil {
		return nil, err
	}
	userDTO, expired, err := validateToken(tokenString, refreshKey)
	if expired {
		RemoveToken(tokenString)
	}
	return userDTO, err
}

// creating token
func createToken(payload *user.UserDTO, expMinutes uint, key string) (string, error) {
	birthdate := ""
	if payload.Birthdate != nil {
		birthdate = payload.Birthdate.Format("02.01.2006")
	}
	expirationTime := time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      payload.ID,
		"username":     payload.Username,
		"email":        payload.Email,
		"phone":        payload.Phone,
		"firstname":    payload.Firstname,
		"lastname":     payload.Lastname,
		"birthdate":    birthdate,
		"is_banned":    payload.IsBanned,
		"is_activated": payload.IsActivated,
		"is_deleted":   payload.IsDeleted,
		"exp":          expirationTime,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		logger.GetInstance().Error(err.Error(), "token creating", map[string]interface{}{"payload": payload}, err)
		return "", appErr.InternalServerError("internal server error")
	}
	return tokenString, nil
}

// validating token
func validateToken(tokenString, key string) (*user.UserDTO, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		expired := false
		if errors.Is(err, jwt.ErrTokenExpired) {
			expired = true
		}
		return nil, expired, appErr.Unauthorized("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userDTO := user.UserDTO{
			ID:          uint64(claims["user_id"].(float64)),
			Username:    claims["username"].(string),
			Email:       claims["email"].(string),
			Firstname:   claims["firstname"].(string),
			Lastname:    claims["lastname"].(string),
			IsBanned:    claims["is_banned"].(bool),
			IsActivated: claims["is_activated"].(bool),
			IsDeleted:   claims["is_deleted"].(bool),
		}
		if phone, ok := claims["phone"].(string); ok {
			userDTO.Phone = utils.StringPtr(phone)
		}
		if birthdate, ok := claims["birthdate"]; ok {
			userDTO.Birthdate = parseDateFromToken(birthdate)
		}
		return &userDTO, false, nil
	}
	return nil, false, appErr.Unauthorized("unauthorized")
}

// parsing date from encrypted user dto object
func parseDateFromToken(date interface{}) *time.Time {
	if date == "" {
		return nil
	}
	parsedDate, err := utils.ParseDate(date.(string))
	if err != nil {
		return nil
	}
	return parsedDate
}
