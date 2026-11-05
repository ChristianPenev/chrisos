package query

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type fakeReader struct {
	views map[user.ID]UserView
}

func (f fakeReader) UserByID(_ context.Context, id user.ID) (UserView, error) {
	v, ok := f.views[id]
	if !ok {
		return UserView{}, user.ErrUserNotFound
	}
	return v, nil
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	id := user.NewID()
	handler := NewGetUserHandler(fakeReader{views: map[user.ID]UserView{
		id: {ID: id.String(), Email: "jane@example.com"},
	}})

	t.Run("returns the view", func(t *testing.T) {
		v, err := handler.Handle(ctx, GetUser{ID: id.String()})
		if err != nil {
			t.Fatal(err)
		}
		if v.Email != "jane@example.com" {
			t.Errorf("got %+v", v)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := handler.Handle(ctx, GetUser{ID: "not-a-uuid"})
		if !errors.Is(err, user.ErrInvalidID) {
			t.Fatalf("want ErrInvalidID, got %v", err)
		}
	})

	t.Run("unknown ID", func(t *testing.T) {
		_, err := handler.Handle(ctx, GetUser{ID: user.NewID().String()})
		if !errors.Is(err, user.ErrUserNotFound) {
			t.Fatalf("want ErrUserNotFound, got %v", err)
		}
	})
}
