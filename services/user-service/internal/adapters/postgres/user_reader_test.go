package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func TestUserReader(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	reader := NewUserReader(db)

	jane := newUser(t, "jane@example.com")
	if err := NewUserRepository(db).Create(ctx, jane); err != nil {
		t.Fatal(err)
	}

	t.Run("UserByID returns the view", func(t *testing.T) {
		v, err := reader.UserByID(ctx, jane.ID())
		if err != nil {
			t.Fatal(err)
		}
		if v.ID != jane.ID().String() ||
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
