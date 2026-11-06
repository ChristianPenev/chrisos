package command

import (
	"context"
	"errors"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type fakeOrganizations struct{ created []*organization.Organization }

func (f *fakeOrganizations) Create(_ context.Context, o *organization.Organization) error {
	f.created = append(f.created, o)
	return nil
}

func (f *fakeOrganizations) FindByID(context.Context, organization.ID) (*organization.Organization, error) {
	return nil, organization.ErrOrganizationNotFound
}

type fakeRoles struct{ created []*role.Role }

func (f *fakeRoles) Create(_ context.Context, r *role.Role) error {
	f.created = append(f.created, r)
	return nil
}

func (f *fakeRoles) FindByID(context.Context, role.ID) (*role.Role, error) {
	return nil, role.ErrRoleNotFound
}

type fakeUsers struct {
	created   []*user.User
	createErr error
}

func (f *fakeUsers) Create(_ context.Context, u *user.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, u)
	return nil
}

func (f *fakeUsers) FindByID(context.Context, user.ID) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (f *fakeUsers) FindByEmail(context.Context, user.Email) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

type fakeTransactor struct {
	repos Repositories
	calls int
}

func (f *fakeTransactor) WithinTransaction(_ context.Context, fn func(Repositories) error) error {
	f.calls++
	return fn(f.repos)
}

type fakes struct {
	orgs  *fakeOrganizations
	roles *fakeRoles
	users *fakeUsers
	tx    *fakeTransactor
}

func newFakes() fakes {
	f := fakes{orgs: &fakeOrganizations{}, roles: &fakeRoles{}, users: &fakeUsers{}}
	f.tx = &fakeTransactor{repos: Repositories{Organizations: f.orgs, Roles: f.roles, Users: f.users}}
	return f
}

func TestCreateOrganization(t *testing.T) {
	ctx := context.Background()
	valid := CreateOrganization{Name: "Acme", Email: "Jane@Example.com", FirstName: "Jane", LastName: "Doe"}

	t.Run("creates the organization, owner role and owner user", func(t *testing.T) {
		f := newFakes()
		res, err := NewCreateOrganizationHandler(f.tx).Handle(ctx, valid)
		if err != nil {
			t.Fatal(err)
		}
		if f.tx.calls != 1 {
			t.Fatalf("want 1 transaction, got %d", f.tx.calls)
		}
		if len(f.orgs.created) != 1 || len(f.roles.created) != 1 || len(f.users.created) != 1 {
			t.Fatalf("created orgs=%d roles=%d users=%d, want 1 each",
				len(f.orgs.created), len(f.roles.created), len(f.users.created))
		}

		org, owner, u := f.orgs.created[0], f.roles.created[0], f.users.created[0]
		if res.OrganizationID != org.ID() || res.UserID != u.ID() {
			t.Errorf("result %+v does not match stored ids", res)
		}
		if org.Name() != "Acme" {
			t.Errorf("org name = %q", org.Name())
		}
		if !owner.IsOwner() || owner.OrganizationID() != org.ID() {
			t.Errorf("role owner=%v org=%v, want owner role of %v", owner.IsOwner(), owner.OrganizationID(), org.ID())
		}
		if u.OrganizationID() != org.ID() || u.RoleID() != owner.ID() {
			t.Errorf("user org=%v role=%v, want %v %v", u.OrganizationID(), u.RoleID(), org.ID(), owner.ID())
		}
		if u.Email().String() != "jane@example.com" {
			t.Errorf("email not normalized: %s", u.Email())
		}
	})

	invalid := []struct {
		name    string
		cmd     CreateOrganization
		wantErr error
	}{
		{"invalid organization name", CreateOrganization{Name: " ", Email: "jane@example.com", FirstName: "Jane", LastName: "Doe"}, organization.ErrInvalidName},
		{"invalid email", CreateOrganization{Name: "Acme", Email: "nope", FirstName: "Jane", LastName: "Doe"}, user.ErrInvalidEmail},
		{"invalid first name", CreateOrganization{Name: "Acme", Email: "jane@example.com", FirstName: "", LastName: "Doe"}, user.ErrInvalidFirstName},
	}
	for _, tt := range invalid {
		t.Run(tt.name+" never opens a transaction", func(t *testing.T) {
			f := newFakes()
			_, err := NewCreateOrganizationHandler(f.tx).Handle(ctx, tt.cmd)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("want %v, got %v", tt.wantErr, err)
			}
			if f.tx.calls != 0 {
				t.Fatal("transaction was opened")
			}
		})
	}

	t.Run("repository errors are returned", func(t *testing.T) {
		f := newFakes()
		f.users.createErr = user.ErrEmailAlreadyExists
		_, err := NewCreateOrganizationHandler(f.tx).Handle(ctx, valid)
		if !errors.Is(err, user.ErrEmailAlreadyExists) {
			t.Fatalf("want ErrEmailAlreadyExists, got %v", err)
		}
	})
}
