package query

import (
	"context"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

type ListUsers struct {
	OrganizationID string
	PageSize       int
	PageToken      string
}

type ListUsersHandler struct {
	reader UserReader
}

func NewListUsersHandler(reader UserReader) ListUsersHandler {
	return ListUsersHandler{reader: reader}
}

func (h ListUsersHandler) Handle(ctx context.Context, q ListUsers) (pagination.Page[UserView], error) {
	orgID, err := organization.ParseID(q.OrganizationID)
	if err != nil {
		return pagination.Page[UserView]{}, err
	}
	req, err := pagination.NewRequest(q.PageSize, q.PageToken)
	if err != nil {
		return pagination.Page[UserView]{}, err
	}
	return h.reader.UsersByOrganization(ctx, orgID, req)
}
