package query

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	orgID := organization.NewID()

	t.Run("passes the organization and normalized page request to the reader", func(t *testing.T) {
		reader := &fakeReader{page: pagination.Page[UserView]{
			Items:      []UserView{{ID: "u1"}},
			NextCursor: "next",
		}}

		page, err := NewListUsersHandler(reader).Handle(ctx, ListUsers{OrganizationID: orgID.String(), PageToken: "abc"})
		if err != nil {
			t.Fatal(err)
		}
		if reader.gotOrg != orgID {
			t.Errorf("reader got org %v, want %v", reader.gotOrg, orgID)
		}
		if reader.gotReq != (pagination.Request{Size: pagination.DefaultSize, After: "abc"}) {
			t.Errorf("reader got request %+v", reader.gotReq)
		}
		if len(page.Items) != 1 || page.NextCursor != "next" {
			t.Errorf("got page %+v", page)
		}
	})

	t.Run("invalid organization ID", func(t *testing.T) {
		_, err := NewListUsersHandler(&fakeReader{}).Handle(ctx, ListUsers{OrganizationID: "nope"})
		if !errors.Is(err, organization.ErrInvalidID) {
			t.Fatalf("want ErrInvalidID, got %v", err)
		}
	})

	t.Run("negative page size", func(t *testing.T) {
		_, err := NewListUsersHandler(&fakeReader{}).Handle(ctx, ListUsers{OrganizationID: orgID.String(), PageSize: -1})
		if !errors.Is(err, pagination.ErrInvalidPageSize) {
			t.Fatalf("want ErrInvalidPageSize, got %v", err)
		}
	})
}
