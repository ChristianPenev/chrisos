package query

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/ChristianPenev/chrisos/pkg/batch"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

type fakeRoleReader struct {
	calls  int
	gotIDs []role.ID
}

func (f *fakeRoleReader) RolesByIDs(_ context.Context, ids []role.ID) ([]RoleView, error) {
	f.calls++
	f.gotIDs = ids
	views := make([]RoleView, len(ids))
	for i, id := range ids {
		views[i] = RoleView{ID: id.String()}
	}
	return views, nil
}

func TestGetRoles(t *testing.T) {
	ctx := context.Background()
	a, b := role.NewID(), role.NewID()

	t.Run("passes parsed, unique IDs to the reader", func(t *testing.T) {
		reader := &fakeRoleReader{}
		views, err := NewGetRolesHandler(reader).Handle(ctx, GetRoles{IDs: []string{a.String(), b.String(), a.String()}})
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(reader.gotIDs, []role.ID{a, b}) {
			t.Errorf("reader got %v", reader.gotIDs)
		}
		if len(views) != 2 {
			t.Errorf("got %d views", len(views))
		}
	})

	t.Run("empty request does not query", func(t *testing.T) {
		reader := &fakeRoleReader{}
		views, err := NewGetRolesHandler(reader).Handle(ctx, GetRoles{})
		if err != nil || len(views) != 0 || views == nil {
			t.Fatalf("got %v, %v", views, err)
		}
		if reader.calls != 0 {
			t.Error("reader was called")
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := NewGetRolesHandler(&fakeRoleReader{}).Handle(ctx, GetRoles{IDs: []string{a.String(), "nope"}})
		if !errors.Is(err, role.ErrInvalidID) {
			t.Fatalf("want ErrInvalidID, got %v", err)
		}
	})

	t.Run("too many IDs", func(t *testing.T) {
		ids := make([]string, batch.MaxIDs+1)
		for i := range ids {
			ids[i] = role.NewID().String()
		}
		_, err := NewGetRolesHandler(&fakeRoleReader{}).Handle(ctx, GetRoles{IDs: ids})
		if !errors.Is(err, batch.ErrTooManyIDs) {
			t.Fatalf("want ErrTooManyIDs, got %v", err)
		}
	})
}
