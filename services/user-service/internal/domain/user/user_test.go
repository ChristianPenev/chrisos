package user

import (
	"errors"
	"strings"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "valid", raw: "jane@example.com", want: "jane@example.com"},
		{name: "normalizes case and spaces", raw: "  Jane.Doe@Example.COM ", want: "jane.doe@example.com"},
		{name: "empty", raw: "   ", wantErr: true},
		{name: "no at sign", raw: "jane.example.com", wantErr: true},
		{name: "display name form", raw: "Jane <jane@example.com>", wantErr: true},
		{name: "too long", raw: strings.Repeat("a", 250) + "@x.com", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.raw)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidEmail) {
					t.Fatalf("want ErrInvalidEmail, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	email, _ := NewEmail("jane@example.com")
	orgID := organization.NewID()
	roleID := role.NewID()

	t.Run("valid", func(t *testing.T) {
		u, err := NewUser(orgID, roleID, email, "  Jane ", "Doe")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.ID() == (ID{}) {
			t.Error("ID not generated")
		}
		if u.OrganizationID() != orgID {
			t.Errorf("OrganizationID() = %v, want %v", u.OrganizationID(), orgID)
		}
		if u.RoleID() != roleID {
			t.Errorf("RoleID() = %v, want %v", u.RoleID(), roleID)
		}
		if u.FirstName() != "Jane" {
			t.Errorf("first name not trimmed: %q", u.FirstName())
		}
		if u.CreatedAt().IsZero() {
			t.Error("createdAt not set")
		}
	})

	tests := []struct {
		name      string
		orgID     organization.ID
		roleID    role.ID
		email     Email
		firstName string
		lastName  string
		wantErr   error
	}{
		{"zero organization", organization.ID{}, roleID, email, "Jane", "Doe", organization.ErrInvalidID},
		{"zero role", orgID, role.ID{}, email, "Jane", "Doe", role.ErrInvalidID},
		{"empty first name", orgID, roleID, email, " ", "Doe", ErrInvalidFirstName},
		{"empty last name", orgID, roleID, email, "Jane", "", ErrInvalidLastName},
		{"too long first name", orgID, roleID, email, strings.Repeat("ж", 101), "Doe", ErrInvalidFirstName},
		{"zero email", orgID, roleID, Email{}, "Jane", "Doe", ErrInvalidEmail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.orgID, tt.roleID, tt.email, tt.firstName, tt.lastName)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("want %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestChangeRole(t *testing.T) {
	email, _ := NewEmail("jane@example.com")
	u, err := NewUser(organization.NewID(), role.NewID(), email, "Jane", "Doe")
	if err != nil {
		t.Fatal(err)
	}

	newRole := role.NewID()
	if err := u.ChangeRole(newRole); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.RoleID() != newRole {
		t.Errorf("RoleID() = %v, want %v", u.RoleID(), newRole)
	}

	if err := u.ChangeRole(role.ID{}); !errors.Is(err, role.ErrInvalidID) {
		t.Errorf("err = %v, want %v", err, role.ErrInvalidID)
	}
	if u.RoleID() != newRole {
		t.Errorf("failed change replaced role with %v", u.RoleID())
	}
}

func TestParseID(t *testing.T) {
	id := NewID()
	parsed, err := ParseID(id.String())
	if err != nil || parsed != id {
		t.Fatalf("round trip failed: %v", err)
	}
	if _, err := ParseID("not-a-uuid"); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("want ErrInvalidID, got %v", err)
	}
}
