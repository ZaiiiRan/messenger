package user

import (
	"strings"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
	privacysettings "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/privacy_settings"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
)

const (
	EmailUpdateCooldown = 24 * time.Hour
)

type User struct {
	id              string
	username        string
	email           string
	profile         *profile.Profile
	status          *status.Status
	privacySettings *privacysettings.PrivacySettings
	createdAt       time.Time
	updatedAt       time.Time
}

func New(
	username, email string,
	profile *profile.Profile,
	privacySettings *privacysettings.PrivacySettings,
	status *status.Status,
) (*User, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)
	u := &User{
		profile:         profile,
		status:          status,
		privacySettings: privacySettings,
	}

	now := time.Now()

	if err := u.SetUsername(username); err != nil {
		verr["username"] = err.Error()
	}
	if err := u.SetEmail(email, &now); err != nil {
		verr["email"] = err.Error()
	}

	if len(verr) > 0 {
		return nil, verr
	}

	u.createdAt = now
	u.updatedAt = now
	return u, nil
}

func FromStorage(
	id, username, email string,
	profile *profile.Profile,
	status *status.Status,
	privacySettings *privacysettings.PrivacySettings,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:              id,
		username:        username,
		email:           email,
		profile:         profile,
		status:          status,
		privacySettings: privacySettings,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

func (u *User) GetID() string                                        { return u.id }
func (u *User) GetUsername() string                                  { return u.username }
func (u *User) GetEmail() string                                     { return u.email }
func (u *User) GetProfile() *profile.Profile                         { return u.profile }
func (u *User) GetPrivacySettings() *privacysettings.PrivacySettings { return u.privacySettings }
func (u *User) GetStatus() *status.Status                            { return u.status }
func (u *User) GetCreatedAt() time.Time                              { return u.createdAt }
func (u *User) GetUpdatedAt() time.Time                              { return u.updatedAt }

func (u *User) SetID(id string) {
	if u.id == "" {
		u.id = id
	}
}

func (u *User) SetUsername(username string) error {
	if u.username == username {
		return ErrSameUsername
	}

	if err := validateUsername(username); err != nil {
		return err
	}
	u.username = strings.ToLower(username)
	return nil
}

func (u *User) SetEmail(email string, now *time.Time) error {
	if u.email == email {
		return ErrSameEmail
	}

	if err := validateEmail(email); err != nil {
		return err
	}

	if time.Since(u.status.GetEmailUpdatedAt()) < EmailUpdateCooldown {
		return ErrWaitBeforeEmailChanging
	}

	var oldEmailPtr *string
	if u.email != "" {
		oldEmail := u.email
		oldEmailPtr = &oldEmail
	}
	u.status.SetOldEmail(oldEmailPtr)
	u.email = strings.ToLower(email)

	if now == nil {
		now = utils.TimePtr(time.Now())
	}
	u.status.SetEmailUpdatedAt(*now)
	return nil
}

func (u *User) SetUpdatedAt(updatedAt *time.Time) {
	if updatedAt == nil {
		u.updatedAt = time.Now()
	} else {
		u.updatedAt = *updatedAt
	}
}
