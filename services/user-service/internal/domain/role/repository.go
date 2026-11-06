package role

import "context"

type Repository interface {
	Create(ctx context.Context, r *Role) error
	FindByID(ctx context.Context, id ID) (*Role, error)
}
