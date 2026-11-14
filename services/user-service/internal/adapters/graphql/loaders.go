package graphql

import (
	"context"
	"errors"

	"github.com/ChristianPenev/chrisos/pkg/dataloader"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
)

type Loaders struct {
	Users *dataloader.Loader[string, query.UserView]
	Roles *dataloader.Loader[string, query.RoleView]
}

func newLoaders(app application.Application) *Loaders {
	return &Loaders{
		Users: dataloader.New(
			func(ctx context.Context, ids []string) ([]query.UserView, error) {
				return app.Queries.GetUsers.Handle(ctx, query.GetUsers{IDs: ids})
			},
			func(v query.UserView) string { return v.ID },
		),
		Roles: dataloader.New(
			func(ctx context.Context, ids []string) ([]query.RoleView, error) {
				return app.Queries.GetRoles.Handle(ctx, query.GetRoles{IDs: ids})
			},
			func(v query.RoleView) string { return v.ID },
		),
	}
}

func loadEntities[V, M any](ctx context.Context, loader *dataloader.Loader[string, V], ids []string, to func(V) *M) ([]*M, error) {
	thunks := make([]func() (V, error), len(ids))
	for i, id := range ids {
		thunks[i] = loader.LoadThunk(ctx, id)
	}

	entities := make([]*M, len(ids))
	for i, thunk := range thunks {
		v, err := thunk()
		if errors.Is(err, dataloader.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		entities[i] = to(v)
	}
	return entities, nil
}
