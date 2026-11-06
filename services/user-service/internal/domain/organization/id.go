package organization

import (
	"fmt"

	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.Must(uuid.NewV7()))
}

func ParseID(s string) (ID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return ID{}, fmt.Errorf("%w: %q", ErrInvalidID, s)
	}
	return ID(u), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}
