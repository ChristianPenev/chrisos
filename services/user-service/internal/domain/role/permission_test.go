package role

import (
	"errors"
	"testing"
)

func TestParsePermission(t *testing.T) {
	for _, p := range AllPermissions() {
		got, err := ParsePermission(p.String())
		if err != nil || got != p {
			t.Errorf("ParsePermission(%q) = %q, %v", p, got, err)
		}
	}

	for _, s := range []string{"", "members.delete", "MEMBERS.INVITE", " members.invite"} {
		if _, err := ParsePermission(s); !errors.Is(err, ErrUnknownPermission) {
			t.Errorf("ParsePermission(%q) err = %v, want %v", s, err, ErrUnknownPermission)
		}
	}
}

func TestAllPermissionsReturnsCopy(t *testing.T) {
	perms := AllPermissions()
	perms[0] = "tampered"

	if AllPermissions()[0] == "tampered" {
		t.Fatal("AllPermissions exposed the internal slice")
	}
}

func TestAllPermissionsUnique(t *testing.T) {
	seen := map[Permission]bool{}
	for _, p := range AllPermissions() {
		if seen[p] {
			t.Errorf("duplicate permission %q", p)
		}
		seen[p] = true
	}
}
