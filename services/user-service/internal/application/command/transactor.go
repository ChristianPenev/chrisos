package command

import (
	"context"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type Repositories struct {
	Organizations organization.Repository
	Roles         role.Repository
	Users         user.Repository
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(Repositories) error) error
}
