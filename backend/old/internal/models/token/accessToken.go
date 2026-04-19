package token

import (
	"backend/internal/models/user/userDTO"
	"os"
)

var accessKey = os.Getenv("JWT_ACCESS_KEY")

// Generate access token
func GenerateAccessToken(payload *userDTO.UserDTO) (string, error) {
	// 30 minutes
	accessToken, err := createToken(payload, 30, accessKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Validating access token
func ValidateAccessToken(tokenString string) (*userDTO.UserDTO, error) {
	userDTO, _, err := validateToken(tokenString, accessKey)
	return userDTO, err
}
