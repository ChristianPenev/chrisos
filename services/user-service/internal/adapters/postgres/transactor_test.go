package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func TestTransactor(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	tx := NewTransactor(db)

	t.Run("commits all changes", func(t *testing.T) {
		org, _ := organization.NewOrganization("Acme")
		owner := role.NewOwnerRole(org.ID())

		err := tx.WithinTransaction(ctx, func(repos command.Repositories) error {
			if err := repos.Organizations.Create(ctx, org); err != nil {
				return err
			}
			return repos.Roles.Create(ctx, owner)
		})
		if err != nil {
			t.Fatal(err)
		}

		if _, err := NewOrganizationRepository(db).FindByID(ctx, org.ID()); err != nil {
			t.Errorf("organization not committed: %v", err)
		}
		if _, err := NewRoleRepository(db).FindByID(ctx, owner.ID()); err != nil {
			t.Errorf("role not committed: %v", err)
		}
	})

	t.Run("rolls back all changes when one fails", func(t *testing.T) {
		orgID, roleID := seedOrganization(t, db)
		jane := newUser(t, orgID, roleID, "taken@example.com")
		if err := NewUserRepository(db).Create(ctx, jane); err != nil {
			t.Fatal(err)
		}

		org, _ := organization.NewOrganization("Globex")
		owner := role.NewOwnerRole(org.ID())
		email, _ := user.NewEmail("taken@example.com")
		dup, _ := user.NewUser(org.ID(), owner.ID(), email, "Jane", "Doe")

		err := tx.WithinTransaction(ctx, func(repos command.Repositories) error {
			if err := repos.Organizations.Create(ctx, org); err != nil {
				return err
			}
			if err := repos.Roles.Create(ctx, owner); err != nil {
				return err
			}
			return repos.Users.Create(ctx, dup)
		})
		if !errors.Is(err, user.ErrEmailAlreadyExists) {
			t.Fatalf("want ErrEmailAlreadyExists, got %v", err)
		}

		if _, err := NewOrganizationRepository(db).FindByID(ctx, org.ID()); !errors.Is(err, organization.ErrOrganizationNotFound) {
			t.Errorf("organization was not rolled back: %v", err)
		}
		if _, err := NewRoleRepository(db).FindByID(ctx, owner.ID()); !errors.Is(err, role.ErrRoleNotFound) {
			t.Errorf("role was not rolled back: %v", err)
		}
	})
}
