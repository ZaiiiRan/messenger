package code

import (
	"time"
)

const (
	maxGenerationsLeft   = 3
	maxVerificationsLeft = 10
	codeTTL              = 10 * time.Minute
	blockWindow          = 5 * time.Minute
	resendCooldown       = 1 * time.Minute
)

type Code struct {
	id                int64
	userId            string
	code              string
	generationsLeft   int
	verificationsLeft int
	expiresAt         time.Time
	createdAt         time.Time
	updatedAt         time.Time
}

func New(userId string) (*Code, error) {
	c := &Code{}
	c.generationsLeft = maxGenerationsLeft
	c.verificationsLeft = maxVerificationsLeft
	c.userId = userId

	if err := c.GenerateCode(); err != nil {
		return nil, err
	}

	now := time.Now()
	c.createdAt = now
	c.updatedAt = now

	return c, nil
}

func FromStorage(
	id int64,
	userId, code string,
	generationsLeft, verificationsLeft int,
	expiresAt, createdAt, updatedAt time.Time,
) *Code {
	return &Code{
		id:                id,
		userId:            userId,
		code:              code,
		generationsLeft:   generationsLeft,
		verificationsLeft: verificationsLeft,
		expiresAt:         expiresAt,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

func (c *Code) GetID() int64               { return c.id }
func (c *Code) GetUserID() string          { return c.userId }
func (c *Code) GetCode() string            { return c.code }
func (c *Code) GetGenerationsLeft() int    { return c.generationsLeft }
func (c *Code) GetVerificationsLeft() int  { return c.verificationsLeft }
func (c *Code) GetExpiresAt() time.Time    { return c.expiresAt }
func (c *Code) GetCreatedAt() time.Time    { return c.createdAt }
func (c *Code) GetUpdatedAt() time.Time    { return c.updatedAt }

func (c *Code) SetID(id int64) {
	if c.id == 0 {
		c.id = id
	}
}

func (c *Code) GenerateCode() error {
	if c.verificationsLeft <= 0 && time.Since(c.updatedAt) < blockWindow {
		return NewCodeValidationError("too many failed confirmation attempts, please try again later")
	}
	if c.verificationsLeft <= 0 {
		c.generationsLeft = maxGenerationsLeft
	}

	if c.generationsLeft <= 0 && time.Since(c.updatedAt) < blockWindow {
		return NewCodeValidationError("the number of code resends has been exhausted")
	}
	if c.generationsLeft <= 0 {
		c.generationsLeft = maxGenerationsLeft
	}

	if !c.updatedAt.IsZero() && time.Since(c.updatedAt) < resendCooldown {
		return NewCodeValidationError("please wait before requesting a new code")
	}

	c.generationsLeft--
	c.verificationsLeft = maxVerificationsLeft

	code, err := generateSixDigitCode()
	if err != nil {
		return err
	}
	c.code = code
	c.expiresAt = time.Now().Add(codeTTL)
	c.updatedAt = time.Now()
	return nil
}

func (c *Code) CheckCode(rawCode string) (bool, error) {
	if time.Now().After(c.expiresAt) {
		return false, NewCodeValidationError("code has been expired")
	}

	if c.verificationsLeft <= 0 && time.Since(c.updatedAt) < blockWindow {
		return false, NewCodeValidationError("too many failed confirmation attempts, please try again later")
	}
	if c.verificationsLeft <= 0 {
		c.verificationsLeft = maxVerificationsLeft
	}

	c.verificationsLeft--

	if c.code == rawCode {
		return true, nil
	}

	if c.verificationsLeft <= 0 {
		c.updatedAt = time.Now()
	}

	return false, nil
}
