package graphql

import (
	"context"
	"fmt"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/model"
)

func (r *entityResolver) FindManyRoleByIDs(ctx context.Context, reps []*model.RoleByIDsInput) ([]*model.Role, error) {
	panic(fmt.Errorf("not implemented: FindManyRoleByIDs - findManyRoleByIDs"))
}

func (r *entityResolver) FindManyUserByIDs(ctx context.Context, reps []*model.UserByIDsInput) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented: FindManyUserByIDs - findManyUserByIDs"))
}

func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
