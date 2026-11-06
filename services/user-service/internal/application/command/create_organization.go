package command

import (
	"context"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type CreateOrganization struct {
	Name      string
	Email     string
	FirstName string
	LastName  string
}

type CreateOrganizationResult struct {
	OrganizationID organization.ID
	UserID         user.ID
}

type CreateOrganizationHandler struct {
	tx Transactor
}

func NewCreateOrganizationHandler(tx Transactor) CreateOrganizationHandler {
	return CreateOrganizationHandler{tx: tx}
}

func (h CreateOrganizationHandler) Handle(ctx context.Context, cmd CreateOrganization) (CreateOrganizationResult, error) {
	org, err := organization.NewOrganization(cmd.Name)
	if err != nil {
		return CreateOrganizationResult{}, err
	}

	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return CreateOrganizationResult{}, err
	}

	owner := role.NewOwnerRole(org.ID())

	u, err := user.NewUser(org.ID(), owner.ID(), email, cmd.FirstName, cmd.LastName)
	if err != nil {
		return CreateOrganizationResult{}, err
	}

	err = h.tx.WithinTransaction(ctx, func(repos Repositories) error {
		if err := repos.Organizations.Create(ctx, org); err != nil {
			return err
		}
		if err := repos.Roles.Create(ctx, owner); err != nil {
			return err
		}
		return repos.Users.Create(ctx, u)
	})
	if err != nil {
		return CreateOrganizationResult{}, err
	}

	return CreateOrganizationResult{OrganizationID: org.ID(), UserID: u.ID()}, nil
}
