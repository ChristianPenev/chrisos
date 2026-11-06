package role

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

func TestNewRole(t *testing.T) {
	orgID := organization.NewID()

	tests := []struct {
		name      string
		roleName  string
		perms     []Permission
		wantName  string
		wantPerms []Permission
		wantErr   error
	}{
		{
			name:      "valid",
			roleName:  "Support",
			perms:     []Permission{PermissionMembersInvite},
			wantName:  "Support",
			wantPerms: []Permission{PermissionMembersInvite},
		},
		{
			name:      "no permissions",
			roleName:  "Viewer",
			wantName:  "Viewer",
			wantPerms: []Permission{},
		},
		{
			name:      "trims name",
			roleName:  "  Support  ",
			wantName:  "Support",
			wantPerms: []Permission{},
		},
		{
			name:      "sorts and removes duplicates",
			roleName:  "Support",
			perms:     []Permission{PermissionRolesManage, PermissionMembersInvite, PermissionRolesManage},
			wantName:  "Support",
			wantPerms: []Permission{PermissionMembersInvite, PermissionRolesManage},
		},
		{name: "empty name", roleName: " ", wantErr: ErrInvalidName},
		{name: "name too long", roleName: strings.Repeat("a", maxNameLength+1), wantErr: ErrInvalidName},
		{name: "unknown permission", roleName: "Support", perms: []Permission{"billing.pay"}, wantErr: ErrUnknownPermission},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRole(orgID, tt.roleName, tt.perms)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", r.Name(), tt.wantName)
			}
			if !slices.Equal(r.Permissions(), tt.wantPerms) {
				t.Errorf("Permissions() = %v, want %v", r.Permissions(), tt.wantPerms)
			}
			if r.OrganizationID() != orgID {
				t.Errorf("OrganizationID() = %v, want %v", r.OrganizationID(), orgID)
			}
			if r.IsOwner() {
				t.Error("IsOwner() = true for a custom role")
			}
			if r.ID() == (ID{}) || r.CreatedAt().IsZero() {
				t.Error("ID or CreatedAt not set")
			}
		})
	}
}

func TestRoleHas(t *testing.T) {
	r, err := NewRole(organization.NewID(), "Support", []Permission{PermissionMembersInvite})
	if err != nil {
		t.Fatal(err)
	}

	if !r.Has(PermissionMembersInvite) {
		t.Error("Has(members.invite) = false, want true")
	}
	if r.Has(PermissionRolesManage) {
		t.Error("Has(roles.manage) = true, want false")
	}
}

func TestOwnerRole(t *testing.T) {
	r := NewOwnerRole(organization.NewID())

	if !r.IsOwner() || r.Name() != OwnerRoleName {
		t.Fatalf("got owner=%v name=%q", r.IsOwner(), r.Name())
	}
	for _, p := range AllPermissions() {
		if !r.Has(p) {
			t.Errorf("owner Has(%q) = false", p)
		}
	}
	if !slices.Equal(r.Permissions(), AllPermissions()) {
		t.Errorf("Permissions() = %v, want all", r.Permissions())
	}
	if err := r.Rename("Boss"); !errors.Is(err, ErrOwnerRoleLocked) {
		t.Errorf("Rename err = %v, want %v", err, ErrOwnerRoleLocked)
	}
	if err := r.SetPermissions(nil); !errors.Is(err, ErrOwnerRoleLocked) {
		t.Errorf("SetPermissions err = %v, want %v", err, ErrOwnerRoleLocked)
	}
}

func TestRename(t *testing.T) {
	r, _ := NewRole(organization.NewID(), "Support", nil)

	if err := r.Rename("  Help desk "); err != nil {
		t.Fatal(err)
	}
	if r.Name() != "Help desk" {
		t.Errorf("Name() = %q", r.Name())
	}
	if err := r.Rename(""); !errors.Is(err, ErrInvalidName) {
		t.Errorf("err = %v, want %v", err, ErrInvalidName)
	}
	if r.Name() != "Help desk" {
		t.Errorf("failed rename changed name to %q", r.Name())
	}
}

func TestSetPermissions(t *testing.T) {
	r, _ := NewRole(organization.NewID(), "Support", []Permission{PermissionMembersInvite})

	if err := r.SetPermissions([]Permission{PermissionMembersRemove}); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(r.Permissions(), []Permission{PermissionMembersRemove}) {
		t.Errorf("Permissions() = %v", r.Permissions())
	}
	if err := r.SetPermissions([]Permission{"billing.pay"}); !errors.Is(err, ErrUnknownPermission) {
		t.Errorf("err = %v, want %v", err, ErrUnknownPermission)
	}
	if !slices.Equal(r.Permissions(), []Permission{PermissionMembersRemove}) {
		t.Errorf("failed update changed permissions to %v", r.Permissions())
	}
}

func TestPermissionsReturnsCopy(t *testing.T) {
	r, _ := NewRole(organization.NewID(), "Support", []Permission{PermissionMembersInvite})

	r.Permissions()[0] = PermissionRolesManage

	if r.Has(PermissionRolesManage) {
		t.Fatal("Permissions() exposed the internal slice")
	}
}

func TestRehydrateSkipsUnknownPermissions(t *testing.T) {
	r := Rehydrate(NewID(), organization.NewID(), "Support",
		[]string{"members.invite", "removed.permission", "members.invite"}, false, time.Now())

	if !slices.Equal(r.Permissions(), []Permission{PermissionMembersInvite}) {
		t.Errorf("Permissions() = %v", r.Permissions())
	}
}
