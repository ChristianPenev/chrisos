package query

import (
	"context"
	"time"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type UserView struct {
	ID             string
	OrganizationID string
	RoleID         string
	Email          string
	FirstName      string
	LastName       string
	CreatedAt      time.Time
}

type UserReader interface {
	UserByID(ctx context.Context, id user.ID) (UserView, error)
}

type GetUser struct {
	ID string
}

type GetUserHandler struct {
	reader UserReader
}

func NewGetUserHandler(reader UserReader) GetUserHandler {
	return GetUserHandler{reader: reader}
}

func (h GetUserHandler) Handle(ctx context.Context, q GetUser) (UserView, error) {
	id, err := user.ParseID(q.ID)
	if err != nil {
		return UserView{}, err
	}
	return h.reader.UserByID(ctx, id)
}
