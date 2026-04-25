package jwt

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Id                  string
	Username            string
	Email               string
	IsConfirmed         bool
	IsDeleted           bool
	IsPermanentlyBanned bool
	IsTemporarilyBanned bool
	Version             int
	jwt.RegisteredClaims
}
