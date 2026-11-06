package graphql

import (
	"context"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input app.RegisterUserInput) (*app.User, error) {
	return r.app.RegisterUser(ctx, input)
}

func (r *queryResolver) User(ctx context.Context, id string) (*app.User, error) {
	return r.app.GetUser(ctx, id)
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
