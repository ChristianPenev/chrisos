package postgres

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func TestUserReader(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	reader := NewUserReader(db)

	orgID, roleID := seedOrganization(t, db)
	jane := newUser(t, orgID, roleID, "jane@example.com")
	if err := NewUserRepository(db).Create(ctx, jane); err != nil {
		t.Fatal(err)
	}

	t.Run("UserByID returns the view", func(t *testing.T) {
		v, err := reader.UserByID(ctx, jane.ID())
		if err != nil {
			t.Fatal(err)
		}
		if v.ID != jane.ID().String() ||
			v.OrganizationID != orgID.String() ||
			v.RoleID != roleID.String() ||
			v.Email != jane.Email().String() ||
			v.FirstName != jane.FirstName() ||
			v.LastName != jane.LastName() ||
			!v.CreatedAt.Equal(jane.CreatedAt()) {
			t.Fatalf("got %+v", v)
		}
	})

	t.Run("UserByID unknown", func(t *testing.T) {
		_, err := reader.UserByID(ctx, user.NewID())
		if !errors.Is(err, user.ErrUserNotFound) {
			t.Fatalf("want ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserReaderUsersByOrganization(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	reader := NewUserReader(db)
	repo := NewUserRepository(db)

	acmeID, acmeRole := seedOrganization(t, db)
	otherID, otherRole := seedOrganization(t, db)

	var want []string
	for i := range 5 {
		u := newUser(t, acmeID, acmeRole, fmt.Sprintf("user%d@acme.com", i))
		if err := repo.Create(ctx, u); err != nil {
			t.Fatal(err)
		}
		want = append(want, u.ID().String())
	}
	if err := repo.Create(ctx, newUser(t, otherID, otherRole, "someone@other.com")); err != nil {
		t.Fatal(err)
	}

	t.Run("pages through the organization's users in creation order", func(t *testing.T) {
		var got []string
		var pages int
		req := pagination.Request{Size: 2}
		for {
			page, err := reader.UsersByOrganization(ctx, acmeID, req)
			if err != nil {
				t.Fatal(err)
			}
			pages++
			for _, v := range page.Items {
				if v.OrganizationID != acmeID.String() {
					t.Fatalf("user %s from organization %s leaked into the list", v.ID, v.OrganizationID)
				}
				got = append(got, v.ID)
			}
			if !page.HasNextPage() {
				break
			}
			req.After = page.NextCursor
		}

		if pages != 3 {
			t.Errorf("pages = %d, want 3", pages)
		}
		if !slices.Equal(got, want) {
			t.Fatalf("got  %v\nwant %v", got, want)
		}
	})

	t.Run("maps all view fields", func(t *testing.T) {
		page, err := reader.UsersByOrganization(ctx, otherID, pagination.Request{Size: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 1 {
			t.Fatalf("got %d users", len(page.Items))
		}
		v := page.Items[0]
		if v.Email != "someone@other.com" || v.RoleID != otherRole.String() || v.FirstName != "Jane" || v.CreatedAt.IsZero() {
			t.Fatalf("got %+v", v)
		}
	})

	t.Run("unknown organization is an empty page", func(t *testing.T) {
		page, err := reader.UsersByOrganization(ctx, organization.NewID(), pagination.Request{Size: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 0 || page.HasNextPage() {
			t.Fatalf("got %+v", page)
		}
	})

	t.Run("invalid cursor", func(t *testing.T) {
		_, err := reader.UsersByOrganization(ctx, acmeID, pagination.Request{Size: 2, After: "garbage"})
		if !errors.Is(err, pagination.ErrInvalidCursor) {
			t.Fatalf("want ErrInvalidCursor, got %v", err)
		}
	})
}
