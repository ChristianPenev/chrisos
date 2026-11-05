package user

import "errors"

var (
	ErrInvalidID        = errors.New("invalid user id")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidFirstName = errors.New("invalid first name")
	ErrInvalidLastName  = errors.New("invalid last name")
)
