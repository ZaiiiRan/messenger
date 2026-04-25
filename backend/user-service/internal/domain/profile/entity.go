package profile

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
)

type Profile struct {
	firstName string
	lastName  string
	phone     *string
	birthdate *time.Time
	bio       *string
}

func New(
	firstName, lastName string,
	phone *string,
	birthdate *time.Time,
	bio *string,
) (*Profile, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)
	p := &Profile{}
	if err := p.SetFirstName(firstName); err != nil {
		verr["profile.first_name"] = err.Error()
	}
	if err := p.SetLastName(lastName); err != nil {
		verr["profile.last_name"] = err.Error()
	}
	if err := p.SetPhone(phone); err != nil {
		verr["profile.phone"] = err.Error()
	}
	if err := p.SetBirthdate(birthdate); err != nil {
		verr["profile.birthdate"] = err.Error()
	}
	p.SetBio(bio)

	if len(verr) > 0 {
		return nil, verr
	}

	return p, nil
}

func FromStorage(
	firstName, lastName string,
	phone *string,
	birthdate *time.Time,
	bio *string,
) *Profile {
	return &Profile{
		firstName: firstName,
		lastName:  lastName,
		phone:     phone,
		birthdate: birthdate,
		bio:       bio,
	}
}

func (p *Profile) GetFirstName() string     { return p.firstName }
func (p *Profile) GetLastName() string      { return p.lastName }
func (p *Profile) GetPhone() *string        { return p.phone }
func (p *Profile) GetBirthdate() *time.Time { return p.birthdate }
func (p *Profile) GetBio() *string          { return p.bio }

func (p *Profile) SetFirstName(firstName string) error {
	if err := validateName(firstName, "first"); err != nil {
		return err
	}
	p.firstName = firstName
	return nil
}

func (p *Profile) SetLastName(lastName string) error {
	if err := validateName(lastName, "last"); err != nil {
		return err
	}
	p.lastName = lastName
	return nil
}

func (p *Profile) SetPhone(phone *string) error {
	if phone == nil || *phone == "" {
		return nil
	}
	if err := validatePhone(*phone); err != nil {
		return err
	}
	p.phone = phone
	return nil
}

func (p *Profile) SetBirthdate(birthdate *time.Time) error {
	if birthdate == nil || birthdate.IsZero() {
		return nil
	}
	if err := validateBirthdate(*birthdate); err != nil {
		return err
	}
	p.birthdate = birthdate
	return nil
}

func (p *Profile) SetBio(bio *string) {
	p.bio = bio
}
