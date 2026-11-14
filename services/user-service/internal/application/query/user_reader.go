package query

import (
	"context"
	"time"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type UserView struct {
	ID             string
	OrganizationID string
	RoleID         string
	Email          string
	FirstName      string
	LastName       string
	CreatedAt      time.Time
}

type UserReader interface {
	UserByID(ctx context.Context, id user.ID) (UserView, error)
	UsersByIDs(ctx context.Context, ids []user.ID) ([]UserView, error)
	UsersByOrganization(ctx context.Context, organizationID organization.ID, req pagination.Request) (pagination.Page[UserView], error)
}
