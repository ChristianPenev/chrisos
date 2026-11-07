package query

import (
	"context"

	"github.com/ChristianPenev/chrisos/pkg/batch"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

type GetRoles struct {
	IDs []string
}

type GetRolesHandler struct {
	reader RoleReader
}

func NewGetRolesHandler(reader RoleReader) GetRolesHandler {
	return GetRolesHandler{reader: reader}
}

func (h GetRolesHandler) Handle(ctx context.Context, q GetRoles) ([]RoleView, error) {
	ids, err := batch.ParseIDs(q.IDs, role.ParseID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []RoleView{}, nil
	}
	return h.reader.RolesByIDs(ctx, ids)
}
