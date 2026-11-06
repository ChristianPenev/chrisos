package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

func TestOrganizationRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewOrganizationRepository(newTestDB(t))

	acme, err := organization.NewOrganization("Acme")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, acme); err != nil {
		t.Fatalf("Create: %v", err)
	}

	t.Run("FindByID returns the same organization", func(t *testing.T) {
		got, err := repo.FindByID(ctx, acme.ID())
		if err != nil {
			t.Fatal(err)
		}
		if got.ID() != acme.ID() || got.Name() != acme.Name() || !got.CreatedAt().Equal(acme.CreatedAt()) {
			t.Fatalf("round trip mismatch:\n got  %+v\n want %+v", *got, *acme)
		}
	})

	t.Run("FindByID unknown", func(t *testing.T) {
		_, err := repo.FindByID(ctx, organization.NewID())
		if !errors.Is(err, organization.ErrOrganizationNotFound) {
			t.Fatalf("want ErrOrganizationNotFound, got %v", err)
		}
	})
}
