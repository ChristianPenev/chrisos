package postgres

import (
	"context"
	"slices"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

func TestRoleReaderRolesByIDs(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	reader := NewRoleReader(db)

	orgID, ownerID := seedOrganization(t, db)
	support, err := role.NewRole(orgID, "Support", []role.Permission{role.PermissionMembersRemove, role.PermissionMembersInvite})
	if err != nil {
		t.Fatal(err)
	}
	if err := NewRoleRepository(db).Create(ctx, support); err != nil {
		t.Fatal(err)
	}

	views, err := reader.RolesByIDs(ctx, []role.ID{support.ID(), ownerID, role.NewID()})
	if err != nil {
		t.Fatal(err)
	}

	byID := map[string]query.RoleView{}
	for _, v := range views {
		byID[v.ID] = v
	}
	if len(views) != 2 {
		t.Fatalf("got %d views, want 2 (unknown ID left out)", len(views))
	}

	t.Run("custom role", func(t *testing.T) {
		v := byID[support.ID().String()]
		if v.Name != "Support" || v.Owner || v.OrganizationID != orgID.String() || v.CreatedAt.IsZero() {
			t.Fatalf("got %+v", v)
		}
		if !slices.Equal(v.Permissions, []string{"members.invite", "members.remove"}) {
			t.Errorf("permissions = %v", v.Permissions)
		}
	})

	t.Run("owner role has every permission", func(t *testing.T) {
		v := byID[ownerID.String()]
		if !v.Owner || v.Name != role.OwnerRoleName {
			t.Fatalf("got %+v", v)
		}
		if len(v.Permissions) != len(role.AllPermissions()) {
			t.Errorf("permissions = %v, want all", v.Permissions)
		}
	})
}
