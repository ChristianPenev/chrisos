package app

import "context"

type Application struct {
	users UserService
}

func New(users UserService) *Application {
	return &Application{users: users}
}

func (a *Application) RegisterUser(ctx context.Context, in RegisterUserInput) (*User, error) {
	id, err := a.users.RegisterUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return a.users.GetUser(ctx, id)
}

func (a *Application) GetUser(ctx context.Context, id string) (*User, error) {
	return a.users.GetUser(ctx, id)
}
