package organization

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const maxNameLength = 100

type Organization struct {
	id        ID
	name      string
	createdAt time.Time
}

func NewOrganization(name string) (*Organization, error) {
	n, err := normalizeName(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidName, err)
	}

	return &Organization{
		id:        NewID(),
		name:      n,
		createdAt: time.Now().UTC().Truncate(time.Microsecond),
	}, nil
}

func Rehydrate(id ID, name string, createdAt time.Time) *Organization {
	return &Organization{
		id:        id,
		name:      name,
		createdAt: createdAt.UTC(),
	}
}

func (o *Organization) ID() ID               { return o.id }
func (o *Organization) Name() string         { return o.name }
func (o *Organization) CreatedAt() time.Time { return o.createdAt }

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
