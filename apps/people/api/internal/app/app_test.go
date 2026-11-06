package app

import (
	"context"
	"errors"
	"testing"
)

type fakeUserService struct {
	users       map[string]*User
	registerErr error
	calls       []string
}

func (f *fakeUserService) RegisterUser(_ context.Context, in RegisterUserInput) (string, error) {
	f.calls = append(f.calls, "RegisterUser")
	if f.registerErr != nil {
		return "", f.registerErr
	}
	f.users["user-1"] = &User{ID: "user-1", Email: in.Email, FirstName: in.FirstName, LastName: in.LastName}
	return "user-1", nil
}

func (f *fakeUserService) GetUser(_ context.Context, id string) (*User, error) {
	f.calls = append(f.calls, "GetUser")
	u, ok := f.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func TestRegisterUserReturnsTheCreatedUser(t *testing.T) {
	fake := &fakeUserService{users: map[string]*User{}}

	u, err := New(fake).RegisterUser(context.Background(), RegisterUserInput{
		Email: "jane@example.com", FirstName: "Jane", LastName: "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != "user-1" || u.Email != "jane@example.com" {
		t.Errorf("got %+v", u)
	}
	if len(fake.calls) != 2 || fake.calls[0] != "RegisterUser" || fake.calls[1] != "GetUser" {
		t.Errorf("unexpected calls: %v", fake.calls)
	}
}

func TestRegisterUserStopsOnError(t *testing.T) {
	fake := &fakeUserService{users: map[string]*User{}, registerErr: ErrEmailAlreadyExists}

	_, err := New(fake).RegisterUser(context.Background(), RegisterUserInput{Email: "jane@example.com"})
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("want ErrEmailAlreadyExists, got %v", err)
	}
	if len(fake.calls) != 1 {
		t.Errorf("GetUser should not be called after a failed register: %v", fake.calls)
	}
}
