package app

import "context"

type Application struct {
	users UserService
}

func New(users UserService) *Application {
	return &Application{users: users}
}

type CreateOrganizationPayload struct {
	OrganizationID string
	Owner          *User
}

func (a *Application) CreateOrganization(ctx context.Context, in CreateOrganizationInput) (*CreateOrganizationPayload, error) {
	created, err := a.users.CreateOrganization(ctx, in)
	if err != nil {
		return nil, err
	}
	owner, err := a.users.GetUser(ctx, created.OwnerUserID)
	if err != nil {
		return nil, err
	}
	return &CreateOrganizationPayload{OrganizationID: created.OrganizationID, Owner: owner}, nil
}

func (a *Application) GetUser(ctx context.Context, id string) (*User, error) {
	return a.users.GetUser(ctx, id)
}

func (a *Application) ListUsers(ctx context.Context, in ListUsersInput) (*UserPage, error) {
	return a.users.ListUsers(ctx, in)
}
