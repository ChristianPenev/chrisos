package postgres

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

func TestRoleRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	repo := NewRoleRepository(db)
	orgID, ownerID := seedOrganization(t, db)

	support, err := role.NewRole(orgID, "Support", []role.Permission{role.PermissionMembersInvite, role.PermissionMembersRemove})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, support); err != nil {
		t.Fatalf("Create: %v", err)
	}

	t.Run("FindByID returns the same role with permissions", func(t *testing.T) {
		got, err := repo.FindByID(ctx, support.ID())
		if err != nil {
			t.Fatal(err)
		}
		assertSameRole(t, got, support)
	})

	t.Run("FindByID returns the owner role", func(t *testing.T) {
		got, err := repo.FindByID(ctx, ownerID)
		if err != nil {
			t.Fatal(err)
		}
		if !got.IsOwner() || got.Name() != role.OwnerRoleName {
			t.Fatalf("got owner=%v name=%q", got.IsOwner(), got.Name())
		}
		if !slices.Equal(got.Permissions(), role.AllPermissions()) {
			t.Errorf("Permissions() = %v, want all", got.Permissions())
		}
	})

	t.Run("role without permissions", func(t *testing.T) {
		viewer, err := role.NewRole(orgID, "Viewer", nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := repo.Create(ctx, viewer); err != nil {
			t.Fatal(err)
		}
		got, err := repo.FindByID(ctx, viewer.ID())
		if err != nil {
			t.Fatal(err)
		}
		assertSameRole(t, got, viewer)
	})

	t.Run("Create with taken name in same organization", func(t *testing.T) {
		dup, _ := role.NewRole(orgID, "Support", nil)
		if err := repo.Create(ctx, dup); !errors.Is(err, role.ErrRoleNameTaken) {
			t.Fatalf("want ErrRoleNameTaken, got %v", err)
		}
	})

	t.Run("Create with owner role name", func(t *testing.T) {
		fake, _ := role.NewRole(orgID, role.OwnerRoleName, nil)
		if err := repo.Create(ctx, fake); !errors.Is(err, role.ErrRoleNameTaken) {
			t.Fatalf("want ErrRoleNameTaken, got %v", err)
		}
	})

	t.Run("FindByID unknown", func(t *testing.T) {
		_, err := repo.FindByID(ctx, role.NewID())
		if !errors.Is(err, role.ErrRoleNotFound) {
			t.Fatalf("want ErrRoleNotFound, got %v", err)
		}
	})

	t.Run("Create in unknown organization", func(t *testing.T) {
		orphan, _ := role.NewRole(organization.NewID(), "Support", nil)
		if err := repo.Create(ctx, orphan); err == nil {
			t.Fatal("want foreign key error, got nil")
		}
	})
}

func assertSameRole(t *testing.T, got, want *role.Role) {
	t.Helper()
	if got.ID() != want.ID() ||
		got.OrganizationID() != want.OrganizationID() ||
		got.Name() != want.Name() ||
		got.IsOwner() != want.IsOwner() ||
		!slices.Equal(got.Permissions(), want.Permissions()) ||
		!got.CreatedAt().Equal(want.CreatedAt()) {
		t.Fatalf("round trip mismatch:\n got  %+v\n want %+v", *got, *want)
	}
}
