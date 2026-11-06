package app

import (
	"context"
	"errors"
	"testing"
)

type fakeUserService struct {
	users     map[string]*User
	createErr error
	calls     []string
}

func (f *fakeUserService) CreateOrganization(_ context.Context, in CreateOrganizationInput) (CreatedOrganization, error) {
	f.calls = append(f.calls, "CreateOrganization")
	if f.createErr != nil {
		return CreatedOrganization{}, f.createErr
	}
	f.users["user-1"] = &User{
		ID: "user-1", OrganizationID: "org-1", RoleID: "role-1",
		Email: in.OwnerEmail, FirstName: in.OwnerFirstName, LastName: in.OwnerLastName,
	}
	return CreatedOrganization{OrganizationID: "org-1", OwnerUserID: "user-1"}, nil
}

func (f *fakeUserService) GetUser(_ context.Context, id string) (*User, error) {
	f.calls = append(f.calls, "GetUser")
	u, ok := f.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func TestCreateOrganizationReturnsTheOwner(t *testing.T) {
	fake := &fakeUserService{users: map[string]*User{}}

	p, err := New(fake).CreateOrganization(context.Background(), CreateOrganizationInput{
		Name: "Acme", OwnerEmail: "jane@example.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.OrganizationID != "org-1" || p.Owner.ID != "user-1" || p.Owner.Email != "jane@example.com" {
		t.Errorf("got %+v owner %+v", p, p.Owner)
	}
	if len(fake.calls) != 2 || fake.calls[0] != "CreateOrganization" || fake.calls[1] != "GetUser" {
		t.Errorf("unexpected calls: %v", fake.calls)
	}
}

func TestCreateOrganizationStopsOnError(t *testing.T) {
	fake := &fakeUserService{users: map[string]*User{}, createErr: ErrEmailAlreadyExists}

	_, err := New(fake).CreateOrganization(context.Background(), CreateOrganizationInput{OwnerEmail: "jane@example.com"})
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("want ErrEmailAlreadyExists, got %v", err)
	}
	if len(fake.calls) != 1 {
		t.Errorf("GetUser should not be called after a failed create: %v", fake.calls)
	}
}
