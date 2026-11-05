package user

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const maxNameLength = 100

type User struct {
	id        ID
	email     Email
	firstName string
	lastName  string
	createdAt time.Time
}


func NewUser(email Email, firstName, lastName string) (*User, error) {
	first, err := normalizeName(firstName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFirstName, err)
	}
	last, err := normalizeName(lastName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidLastName, err)
	}
	if email == (Email{}) {
		return nil, fmt.Errorf("%w: must not be empty", ErrInvalidEmail)
	}

	return &User{
		id:        NewID(),
		email:     email,
		firstName: first,
		lastName:  last,
		createdAt: time.Now().UTC().Truncate(time.Microsecond),
	}, nil
}

func (u *User) ID() ID               { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) FirstName() string    { return u.firstName }
func (u *User) LastName() string     { return u.lastName }
func (u *User) CreatedAt() time.Time { return u.createdAt }

func normalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", fmt.Errorf("must not be empty")
	}
	if utf8.RuneCountInString(name) > maxNameLength {
		return "", fmt.Errorf("longer than %d characters", maxNameLength)
	}
	return name, nil
}
