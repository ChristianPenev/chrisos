package command

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type fakeRepo struct {
	created   []*user.User
	createErr error
}

func (f *fakeRepo) Create(_ context.Context, u *user.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, u)
	return nil
}

func (f *fakeRepo) FindByID(context.Context, user.ID) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (f *fakeRepo) FindByEmail(context.Context, user.Email) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	valid := RegisterUser{Email: "Jane@Example.com", FirstName: "Jane", LastName: "Doe"}

	t.Run("creates the user and returns its ID", func(t *testing.T) {
		repo := &fakeRepo{}
		id, err := NewRegisterUserHandler(repo).Handle(ctx, valid)
		if err != nil {
			t.Fatal(err)
		}
		if len(repo.created) != 1 {
			t.Fatalf("want 1 created user, got %d", len(repo.created))
		}
		got := repo.created[0]
		if got.ID() != id {
			t.Errorf("returned ID %s, stored ID %s", id, got.ID())
		}
		if got.Email().String() != "jane@example.com" {
			t.Errorf("email not normalized: %s", got.Email())
		}
	})

	t.Run("invalid input never reaches the repository", func(t *testing.T) {
		repo := &fakeRepo{}
		_, err := NewRegisterUserHandler(repo).Handle(ctx, RegisterUser{Email: "nope", FirstName: "Jane", LastName: "Doe"})
		if !errors.Is(err, user.ErrInvalidEmail) {
			t.Fatalf("want ErrInvalidEmail, got %v", err)
		}
		if len(repo.created) != 0 {
			t.Fatal("repository was called")
		}
	})

	t.Run("repository errors are returned", func(t *testing.T) {
		repo := &fakeRepo{createErr: user.ErrEmailAlreadyExists}
		_, err := NewRegisterUserHandler(repo).Handle(ctx, valid)
		if !errors.Is(err, user.ErrEmailAlreadyExists) {
			t.Fatalf("want ErrEmailAlreadyExists, got %v", err)
		}
	})
}
