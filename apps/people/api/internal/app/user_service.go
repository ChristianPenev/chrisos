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

type UserService interface {
	CreateOrganization(ctx context.Context, in CreateOrganizationInput) (CreatedOrganization, error)
	GetUser(ctx context.Context, id string) (*User, error)
}
