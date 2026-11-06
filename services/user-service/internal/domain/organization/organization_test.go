package organization

import (
	"errors"
	"strings"
	"testing"
)

func TestNewOrganization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantErr  error
	}{
		{name: "valid", input: "Acme", wantName: "Acme"},
		{name: "trims whitespace", input: "  Acme Inc  ", wantName: "Acme Inc"},
		{name: "max length", input: strings.Repeat("a", maxNameLength), wantName: strings.Repeat("a", maxNameLength)},
		{name: "empty", input: "", wantErr: ErrInvalidName},
		{name: "only whitespace", input: "   ", wantErr: ErrInvalidName},
		{name: "too long", input: strings.Repeat("a", maxNameLength+1), wantErr: ErrInvalidName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org, err := NewOrganization(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if org.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", org.Name(), tt.wantName)
			}
			if org.ID() == (ID{}) {
				t.Error("ID() is zero")
			}
			if org.CreatedAt().IsZero() {
				t.Error("CreatedAt() is zero")
			}
		})
	}
}

func TestParseID(t *testing.T) {
	id := NewID()
	parsed, err := ParseID(id.String())
	if err != nil || parsed != id {
		t.Fatalf("ParseID(%q) = %v, %v", id, parsed, err)
	}

	if _, err := ParseID("not-a-uuid"); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidID)
	}
}
