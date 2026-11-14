package graphql

import (
	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/model"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
)

func toUser(v query.UserView) *model.User {
	return &model.User{
		ID:             v.ID,
		OrganizationID: v.OrganizationID,
		RoleID:         v.RoleID,
		Email:          v.Email,
		FirstName:      v.FirstName,
		LastName:       v.LastName,
		CreatedAt:      v.CreatedAt,
	}
}

func toRole(v query.RoleView) *model.Role {
	return &model.Role{
		ID:             v.ID,
		OrganizationID: v.OrganizationID,
		Name:           v.Name,
		IsOwner:        v.Owner,
		Permissions:    v.Permissions,
		CreatedAt:      v.CreatedAt,
	}
}

func toUserConnection(page pagination.Page[query.UserView]) *model.UserConnection {
	nodes := make([]*model.User, len(page.Items))
	for i, v := range page.Items {
		nodes[i] = toUser(v)
	}

	info := &model.PageInfo{HasNextPage: page.HasNextPage()}
	if info.HasNextPage {
		info.EndCursor = &page.NextCursor
	}
	return &model.UserConnection{Nodes: nodes, PageInfo: info}
}

func valueOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}
