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

type RegisterUserInput struct {
	Email     string
	FirstName string
	LastName  string
}

type UserService interface {
	RegisterUser(ctx context.Context, in RegisterUserInput) (string, error)
	GetUser(ctx context.Context, id string) (*User, error)
}
