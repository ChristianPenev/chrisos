package app

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type CreateOrganizationInput struct {
	Name           string
	OwnerEmail     string
	OwnerFirstName string
	OwnerLastName  string
}

type CreatedOrganization struct {
	OrganizationID string
	OwnerUserID    string
}

type ListUsersInput struct {
	OrganizationID string
	First          int
	After          string
}

type UserPage struct {
	Users      []*User
	NextCursor string
}

type UserService interface {
	CreateOrganization(ctx context.Context, in CreateOrganizationInput) (CreatedOrganization, error)
	GetUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, in ListUsersInput) (*UserPage, error)
	GetRoles(ctx context.Context, ids []string) ([]*Role, error)
}
