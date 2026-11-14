package graphql

import (
	"context"

	"github.com/ChristianPenev/chrisos/pkg/dataloader"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/model"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
)

func (r *mutationResolver) CreateOrganization(ctx context.Context, input model.CreateOrganizationInput) (*model.CreateOrganizationPayload, error) {
	created, err := r.app.Commands.CreateOrganization.Handle(ctx, command.CreateOrganization{
		Name:      input.Name,
		Email:     input.OwnerEmail,
		FirstName: input.OwnerFirstName,
		LastName:  input.OwnerLastName,
	})
	if err != nil {
		return nil, err
	}

	owner, err := r.app.Queries.GetUser.Handle(ctx, query.GetUser{ID: created.UserID.String()})
	if err != nil {
		return nil, err
	}

	return &model.CreateOrganizationPayload{
		OrganizationID: created.OrganizationID.String(),
		Owner:          toUser(owner),
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	v, err := r.app.Queries.GetUser.Handle(ctx, query.GetUser{ID: id})
	if err != nil {
		return nil, err
	}
	return toUser(v), nil
}

func (r *queryResolver) Users(ctx context.Context, organizationID string, first *int, after *string) (*model.UserConnection, error) {
	page, err := r.app.Queries.ListUsers.Handle(ctx, query.ListUsers{
		OrganizationID: organizationID,
		PageSize:       valueOr(first, 0),
		PageToken:      valueOr(after, ""),
	})
	if err != nil {
		return nil, err
	}
	return toUserConnection(page), nil
}

func (r *userResolver) Role(ctx context.Context, obj *model.User) (*model.Role, error) {
	v, err := dataloader.For[*Loaders](ctx).Roles.Load(ctx, obj.RoleID)
	if err != nil {
		return nil, err
	}
	return toRole(v), nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	userResolver     struct{ *Resolver }
)
