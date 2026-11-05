package user

import (
	"errors"
	"strings"
	"testing"
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

	t.Run("valid", func(t *testing.T) {
		u, err := NewUser(email, "  Jane ", "Doe")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.ID() == (ID{}) {
			t.Error("ID not generated")
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
		email     Email
		firstName string
		lastName  string
		wantErr   error
	}{
		{"empty first name", email, " ", "Doe", ErrInvalidFirstName},
		{"empty last name", email, "Jane", "", ErrInvalidLastName},
		{"too long first name", email, strings.Repeat("ж", 101), "Doe", ErrInvalidFirstName},
		{"zero email", Email{}, "Jane", "Doe", ErrInvalidEmail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.email, tt.firstName, tt.lastName)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("want %v, got %v", tt.wantErr, err)
			}
		})
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
