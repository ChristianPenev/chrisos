package graphql

import (
	"context"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/model"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
	"github.com/ChristianPenev/chrisos/pkg/dataloader"
)

func (r *mutationResolver) CreateOrganization(ctx context.Context, input app.CreateOrganizationInput) (*app.CreateOrganizationPayload, error) {
	return r.app.CreateOrganization(ctx, input)
}

func (r *queryResolver) User(ctx context.Context, id string) (*app.User, error) {
	return r.app.GetUser(ctx, id)
}

func (r *queryResolver) Users(ctx context.Context, organizationID string, first *int, after *string) (*model.UserConnection, error) {
	page, err := r.app.ListUsers(ctx, app.ListUsersInput{
		OrganizationID: organizationID,
		First:          valueOr(first, 0),
		After:          valueOr(after, ""),
	})
	if err != nil {
		return nil, err
	}
	return toUserConnection(page), nil
}

func (r *userResolver) Role(ctx context.Context, obj *app.User) (*app.Role, error) {
	return dataloader.For[*Loaders](ctx).Roles.Load(ctx, obj.RoleID)
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	userResolver     struct{ *Resolver }
)
