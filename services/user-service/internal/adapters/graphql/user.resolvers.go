package graphql

import (
	"context"
	"fmt"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/model"
)

func (r *mutationResolver) CreateOrganization(ctx context.Context, input model.CreateOrganizationInput) (*model.CreateOrganizationPayload, error) {
	panic(fmt.Errorf("not implemented: CreateOrganization - createOrganization"))
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

func (r *queryResolver) Users(ctx context.Context, organizationID string, first *int, after *string) (*model.UserConnection, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}

func (r *userResolver) Role(ctx context.Context, obj *model.User) (*model.Role, error) {
	panic(fmt.Errorf("not implemented: Role - role"))
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	userResolver     struct{ *Resolver }
)
