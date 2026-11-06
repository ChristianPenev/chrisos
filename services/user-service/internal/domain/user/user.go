package user

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

const maxNameLength = 100

type User struct {
	id             ID
	organizationID organization.ID
	roleID         role.ID
	email          Email
	firstName      string
	lastName       string
	createdAt      time.Time
}

func NewUser(organizationID organization.ID, roleID role.ID, email Email, firstName, lastName string) (*User, error) {
	if organizationID == (organization.ID{}) {
		return nil, fmt.Errorf("%w: must not be empty", organization.ErrInvalidID)
	}
	if roleID == (role.ID{}) {
		return nil, fmt.Errorf("%w: must not be empty", role.ErrInvalidID)
	}
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
		id:             NewID(),
		organizationID: organizationID,
		roleID:         roleID,
		email:          email,
		firstName:      first,
		lastName:       last,
		createdAt:      time.Now().UTC().Truncate(time.Microsecond),
	}, nil
}

func Rehydrate(id ID, organizationID organization.ID, roleID role.ID, email, firstName, lastName string, createdAt time.Time) *User {
	return &User{
		id:             id,
		organizationID: organizationID,
		roleID:         roleID,
		email:          Email{value: email},
		firstName:      firstName,
		lastName:       lastName,
		createdAt:      createdAt.UTC(),
	}
}

func (u *User) ID() ID                          { return u.id }
func (u *User) OrganizationID() organization.ID { return u.organizationID }
func (u *User) RoleID() role.ID                 { return u.roleID }
func (u *User) Email() Email                    { return u.email }
func (u *User) FirstName() string               { return u.firstName }
func (u *User) LastName() string                { return u.lastName }
func (u *User) CreatedAt() time.Time            { return u.createdAt }

func (u *User) ChangeRole(roleID role.ID) error {
	if roleID == (role.ID{}) {
		return fmt.Errorf("%w: must not be empty", role.ErrInvalidID)
	}
	u.roleID = roleID
	return nil
}

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
