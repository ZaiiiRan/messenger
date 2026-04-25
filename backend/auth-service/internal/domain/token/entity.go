package token

import (
	"time"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type Token struct {
	id        int64
	userId    string
	token     string
	tokenType string
	version   int
	expiresAt time.Time
	createdAt time.Time
	updatedAt time.Time
}

func New(
	userId string,
	token string, tokenType string,
	version int,
	expiresAt time.Time,
) *Token {
	now := time.Now()

	return &Token{
		userId:    userId,
		token:     token,
		tokenType: tokenType,
		version:   version,
		expiresAt: expiresAt,
		createdAt: now,
		updatedAt: now,
	}
}

func FromStorage(
	id int64,
	userId string,
	token string, tokenType string,
	version int,
	expiresAt time.Time, createdAt time.Time, updatedAt time.Time,
) *Token {
	return &Token{
		id:        id,
		userId:    userId,
		token:     token,
		tokenType: tokenType,
		version:   version,
		expiresAt: expiresAt,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t *Token) Id() int64            { return t.id }
func (t *Token) UserId() string       { return t.userId }
func (t *Token) Token() string        { return t.token }
func (t *Token) TokenType() string    { return t.tokenType }
func (t *Token) Version() int         { return t.version }
func (t *Token) ExpiresAt() time.Time { return t.expiresAt }
func (t *Token) CreatedAt() time.Time { return t.createdAt }
func (t *Token) UpdatedAt() time.Time { return t.updatedAt }

func (t *Token) SetId(id int64) {
	if t.Id() == 0 {
		t.id = id
	}
}

func (t *Token) SetToken(token string, expiresAt time.Time) {
	t.token = token
	t.expiresAt = expiresAt
	t.updatedAt = time.Now()
}
