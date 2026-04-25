package jwt

import "github.com/golang-jwt/jwt/v5"

func ParseUserToken(tokenStr string, key []byte) (*UserClaims, error) {
	return parseToken(tokenStr, key, func() *UserClaims { return &UserClaims{} })
}

func parseToken[T jwt.Claims](
	tokenStr string,
	key []byte,
	newClaims func() T,
) (T, error) {
	claims := newClaims()
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			var zero T
			return zero, ErrInvalidToken
		}
		return key, nil
	})

	if err != nil || !token.Valid {
		var zero T
		return zero, ErrInvalidToken
	}

	return claims, nil
}
