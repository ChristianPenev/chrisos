package query

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/ChristianPenev/chrisos/pkg/batch"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func TestGetUsers(t *testing.T) {
	ctx := context.Background()
	a, b := user.NewID(), user.NewID()
	views := map[user.ID]UserView{a: {ID: a.String()}, b: {ID: b.String()}}

	t.Run("passes parsed, unique IDs to the reader", func(t *testing.T) {
		reader := &fakeReader{views: views}
		got, err := NewGetUsersHandler(reader).Handle(ctx, GetUsers{IDs: []string{a.String(), b.String(), a.String()}})
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(reader.gotIDs, []user.ID{a, b}) {
			t.Errorf("reader got %v", reader.gotIDs)
		}
		if len(got) != 2 {
			t.Errorf("got %d views", len(got))
		}
	})

	t.Run("empty request does not query", func(t *testing.T) {
		reader := &fakeReader{}
		got, err := NewGetUsersHandler(reader).Handle(ctx, GetUsers{})
		if err != nil || len(got) != 0 || got == nil {
			t.Fatalf("got %v, %v", got, err)
		}
		if reader.calls != 0 {
			t.Error("reader was called")
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := NewGetUsersHandler(&fakeReader{}).Handle(ctx, GetUsers{IDs: []string{a.String(), "nope"}})
		if !errors.Is(err, user.ErrInvalidID) {
			t.Fatalf("want ErrInvalidID, got %v", err)
		}
	})

	t.Run("too many IDs", func(t *testing.T) {
		ids := make([]string, batch.MaxIDs+1)
		for i := range ids {
			ids[i] = user.NewID().String()
		}
		_, err := NewGetUsersHandler(&fakeReader{}).Handle(ctx, GetUsers{IDs: ids})
		if !errors.Is(err, batch.ErrTooManyIDs) {
			t.Fatalf("want ErrTooManyIDs, got %v", err)
		}
	})
}
