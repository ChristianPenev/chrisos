package role

import "errors"

var (
	ErrInvalidID         = errors.New("invalid role id")
	ErrInvalidName       = errors.New("invalid role name")
	ErrUnknownPermission = errors.New("unknown permission")
	ErrOwnerRoleLocked   = errors.New("owner role cannot be changed")
)
