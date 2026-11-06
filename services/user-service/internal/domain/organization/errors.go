package organization

import "errors"

var (
	ErrInvalidID   = errors.New("invalid organization id")
	ErrInvalidName = errors.New("invalid organization name")
)
