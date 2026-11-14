package query

import (
	"context"

	"github.com/ChristianPenev/chrisos/pkg/batch"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type GetUsers struct {
	IDs []string
}

type GetUsersHandler struct {
	reader UserReader
}

func NewGetUsersHandler(reader UserReader) GetUsersHandler {
	return GetUsersHandler{reader: reader}
}

func (h GetUsersHandler) Handle(ctx context.Context, q GetUsers) ([]UserView, error) {
	ids, err := batch.ParseIDs(q.IDs, user.ParseID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []UserView{}, nil
	}
	return h.reader.UsersByIDs(ctx, ids)
}
