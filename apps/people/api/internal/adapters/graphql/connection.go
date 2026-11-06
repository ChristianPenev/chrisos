package graphql

import (
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/model"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
)

func toUserConnection(page *app.UserPage) *model.UserConnection {
	info := &model.PageInfo{HasNextPage: page.NextCursor != ""}
	if info.HasNextPage {
		info.EndCursor = &page.NextCursor
	}
	return &model.UserConnection{Nodes: page.Users, PageInfo: info}
}

func valueOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}
