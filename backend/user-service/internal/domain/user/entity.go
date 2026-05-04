package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
)

type User struct {
	id        string
	username  string
	email     string
	profile   *profile.Profile
	status    *status.Status
	createdAt time.Time
	updatedAt time.Time
}

func New(
	username, email string,
	profile *profile.Profile,
	status *status.Status,
) (*User, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)
	u := &User{
		profile: profile,
		status:  status,
	}

	if err := u.SetUsername(username); err != nil {
		verr["username"] = err.Error()
	}
	if err := u.SetEmail(email); err != nil {
		verr["email"] = err.Error()
	}

	if len(verr) > 0 {
		return nil, verr
	}

	now := time.Now()
	u.createdAt = now
	u.updatedAt = now
	return u, nil
}

func FromStorage(
	id, username, email string,
	profile *profile.Profile,
	status *status.Status,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:       id,
		username: username,
		email:    email,
		profile:  profile,
		status:   status,
	}
}

func (u *User) GetID() string                { return u.id }
func (u *User) GetUsername() string          { return u.username }
func (u *User) GetEmail() string             { return u.email }
func (u *User) GetProfile() *profile.Profile { return u.profile }
func (u *User) GetStatus() *status.Status    { return u.status }
func (u *User) GetCreatedAt() time.Time      { return u.createdAt }
func (u *User) GetUpdatedAt() time.Time      { return u.updatedAt }

func (u *User) SetID(id string) {
	if u.id == "" {
		u.id = id
	}
}

func (u *User) SetUsername(username string) error {
	if u.username == username {
		return fmt.Errorf("username is the same as the current one")
	}

	if err := validateUsername(username); err != nil {
		return err
	}
	u.username = strings.ToLower(username)
	return nil
}

func (u *User) SetEmail(email string) error {
	if u.email == email {
		return fmt.Errorf("email is the same as the current one")
	}

	if err := validateEmail(email); err != nil {
		return err
	}
	u.email = strings.ToLower(email)
	return nil
}

func (u *User) SetUpdatedAt(updatedAt *time.Time) {
	if updatedAt == nil {
		u.updatedAt = time.Now()
	} else {
		u.updatedAt = *updatedAt
	}
}
