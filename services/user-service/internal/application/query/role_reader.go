package query

import (
	"context"
	"time"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

type RoleView struct {
	ID             string
	OrganizationID string
	Name           string
	Owner          bool
	Permissions    []string
	CreatedAt      time.Time
}

type RoleReader interface {
	RolesByIDs(ctx context.Context, ids []role.ID) ([]RoleView, error)
}
