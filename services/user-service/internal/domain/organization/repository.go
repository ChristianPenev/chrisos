package organization

import "context"

type Repository interface {
	Create(ctx context.Context, o *Organization) error
	FindByID(ctx context.Context, id ID) (*Organization, error)
}
