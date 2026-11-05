package command

import (
	"context"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type RegisterUser struct {
	Email     string
	FirstName string
	LastName  string
}

type RegisterUserHandler struct {
	users user.Repository
}

func NewRegisterUserHandler(users user.Repository) RegisterUserHandler {
	return RegisterUserHandler{users: users}
}

func (h RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUser) (user.ID, error) {
	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return user.ID{}, err
	}

	u, err := user.NewUser(email, cmd.FirstName, cmd.LastName)
	if err != nil {
		return user.ID{}, err
	}

	if err := h.users.Create(ctx, u); err != nil {
		return user.ID{}, err
	}
	return u.ID(), nil
}
