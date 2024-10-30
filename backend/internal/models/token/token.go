package token

type Token struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}
