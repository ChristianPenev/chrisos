package query

import (
	"context"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

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
