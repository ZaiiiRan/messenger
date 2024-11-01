package token

import (
	pgDB "backend/internal/dbs/pgDB"
	dto "backend/internal/dtos/userDTO"
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
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Remove token from DataBase
func RemoveToken(refreshToken string) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`DELETE FROM tokens WHERE refresh_token = $1`, refreshToken)
	return err
}

// Save token in DataBase
func SaveToken(userID uint64, refreshToken string) (*Token, error) {
	db := pgDB.GetDB()
	var token Token
	err := db.QueryRow(`INSERT INTO tokens (user_id, refresh_token) VALUES ($1, $2) RETURNING *`, userID, refreshToken).Scan(
		&token.ID, &token.UserID, &token.RefreshToken)
	return &token, err
}

// Generate token pair
func GenerateTokens(payload *dto.UserDTO) (string, string, error) {
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
func ValidateAccessToken(tokenString string) (*dto.UserDTO, error) {
	userDTO, err := validateToken(tokenString, accessKey)
	return userDTO, err
}

// Validating refresh token
func ValidateRefreshToken(tokenString string) (*dto.UserDTO, error) {
	userDTO, err := validateToken(tokenString, refreshKey)
	return userDTO, err
}

// creating token
func createToken(payload *dto.UserDTO, expMinutes uint, key string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      payload.ID,
		"username":     payload.Username,
		"email":        payload.Email,
		"phone":        payload.Phone,
		"firstname":    payload.Firstname,
		"lastname":     payload.Lastname,
		"birthdate":    payload.Birthdate,
		"is_banned":    payload.IsBanned,
		"is_activated": payload.IsActivated,
		"is_deleted":   payload.IsDeleted,
		"exp":          expirationTime,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

// validating token
func validateToken(tokenString, key string) (*dto.UserDTO, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userDTO := dto.UserDTO{
			ID:          uint64(claims["user_id"].(float64)),
			Username:    claims["usename"].(string),
			Email:       claims["email"].(string),
			Phone:       claims["phone"].(*string),
			Firstname:   claims["firstname"].(string),
			Lastname:    claims["lastname"].(string),
			Birthdate:   claims["birthdate"].(*time.Time),
			IsBanned:    claims["is_banned"].(bool),
			IsActivated: claims["is_activated"].(bool),
			IsDeleted:   claims["is_deleted"].(bool),
		}
		return &userDTO, nil
	}
	return nil, errors.New("unathorized")
}
