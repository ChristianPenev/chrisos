package graphql

import (
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
	"github.com/ChristianPenev/chrisos/pkg/dataloader"
)

type Loaders struct {
	Roles *dataloader.Loader[string, *app.Role]
}

func newLoaders(application *app.Application) *Loaders {
	return &Loaders{
		Roles: dataloader.New(application.GetRoles, func(r *app.Role) string { return r.ID }),
	}
}
