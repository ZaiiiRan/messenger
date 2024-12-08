package token

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/user"
	"backend/internal/models/user/userDTO"
	"backend/internal/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// creating token
func createToken(payload *userDTO.UserDTO, expMinutes uint, key string) (string, error) {
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
func validateToken(tokenString, key string) (*userDTO.UserDTO, bool, error) {
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
		userDTO := userDTO.UserDTO{
			BaseUser: user.BaseUser{
				ID:          uint64(claims["user_id"].(float64)),
				Username:    claims["username"].(string),
				Firstname:   claims["firstname"].(string),
				Lastname:    claims["lastname"].(string),
				IsBanned:    claims["is_banned"].(bool),
				IsActivated: claims["is_activated"].(bool),
				IsDeleted:   claims["is_deleted"].(bool),
			},
			Email:       claims["email"].(string),
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
