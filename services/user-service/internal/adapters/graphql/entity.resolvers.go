package graphql

import (
	"context"

	"github.com/ChristianPenev/chrisos/pkg/dataloader"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/model"
)

func (r *entityResolver) FindManyRoleByIDs(ctx context.Context, reps []*model.RoleByIDsInput) ([]*model.Role, error) {
	ids := make([]string, len(reps))
	for i, rep := range reps {
		ids[i] = rep.ID
	}
	return loadEntities(ctx, dataloader.For[*Loaders](ctx).Roles, ids, toRole)
}

func (r *entityResolver) FindManyUserByIDs(ctx context.Context, reps []*model.UserByIDsInput) ([]*model.User, error) {
	ids := make([]string, len(reps))
	for i, rep := range reps {
		ids[i] = rep.ID
	}
	return loadEntities(ctx, dataloader.For[*Loaders](ctx).Users, ids, toUser)
}

func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
