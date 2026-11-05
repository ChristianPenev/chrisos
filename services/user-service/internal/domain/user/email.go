package user

import (
	"fmt"
	"net/mail"
	"strings"
)

const maxEmailLength = 254

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))

	if normalized == "" {
		return Email{}, fmt.Errorf("%w: must not be empty", ErrInvalidEmail)
	}
	if len(normalized) > maxEmailLength {
		return Email{}, fmt.Errorf("%w: longer than %d characters", ErrInvalidEmail, maxEmailLength)
	}

	addr, err := mail.ParseAddress(normalized)
	if err != nil || addr.Address != normalized {
		return Email{}, fmt.Errorf("%w: %q", ErrInvalidEmail, raw)
	}

	return Email{value: normalized}, nil
}

func (e Email) String() string {
	return e.value
}
