package token

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"os"
)

var refreshKey = os.Getenv("JWT_REFRESH_KEY")

type Token struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// Find refresh token in DataBase
func FindRefreshToken(refreshToken string) (*Token, error) {
	token, err := getTokenFromDB(refreshToken)
	if err != nil && err.Error() == "token not found" {
		return nil, appErr.Unauthorized("unauthorized")
	} else if err != nil {
		return nil, err
	}
	return token, nil
}

// Remove refresh token from DataBase
func (t *Token) RemoveRefreshToken() error {
	return removeTokenFromDB(t)
}

// Remove other tokens (sessions) from DataBase
func (t *Token) RemoveOtherTokens() error {
	return removeOtherTokensFromDB(t)
}

// Save refresh token to DataBase
func (t *Token) SaveRefreshToken() error {
	if t.ID == 0 {
		err := insertTokenToDB(t)
		if err != nil {
			return err
		}
	} else {
		err := updateTokenInDB(t)
		if err != nil {
			return appErr.Unauthorized("unauthorized")
		}
	}
	return nil
}

// Generate refresh token
func GenerateRefreshToken(payload *user.UserDTO) (*Token, error) {
	refreshToken, err := generateRefreshTokenString(payload)
	if err != nil {
		return nil, err
	}

	return &Token{
		UserID:       payload.ID,
		RefreshToken: refreshToken,
	}, nil
}

// Regenerate refresh token
func (t *Token) RegenerateRefreshToken(payload *user.UserDTO) error {
	refreshToken, err := generateRefreshTokenString(payload)
	if err != nil {
		return err
	}

	t.RefreshToken = refreshToken
	return t.SaveRefreshToken()
}

// Validating refresh token
func (t *Token) ValidateRefreshToken() (*user.UserDTO, error) {
	userDTO, expired, err := validateToken(t.RefreshToken, refreshKey)
	if expired {
		t.RemoveRefreshToken()
	}

	return userDTO, err
}

// generate refresh token string
func generateRefreshTokenString(payload *user.UserDTO) (string, error) {
	// 30 days
	refreshToken, err := createToken(payload, 43200, refreshKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}