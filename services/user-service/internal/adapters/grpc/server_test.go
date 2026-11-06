package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type memoryStore struct {
	orgs      map[organization.ID]*organization.Organization
	roles     map[role.ID]*role.Role
	users     map[user.ID]*user.User
	createErr error
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		orgs:  map[organization.ID]*organization.Organization{},
		roles: map[role.ID]*role.Role{},
		users: map[user.ID]*user.User{},
	}
}

func (m *memoryStore) WithinTransaction(_ context.Context, fn func(command.Repositories) error) error {
	return fn(command.Repositories{
		Organizations: memoryOrganizations{m},
		Roles:         memoryRoles{m},
		Users:         memoryUsers{m},
	})
}

func (m *memoryStore) UserByID(_ context.Context, id user.ID) (query.UserView, error) {
	u, ok := m.users[id]
	if !ok {
		return query.UserView{}, user.ErrUserNotFound
	}
	return query.UserView{
		ID:             u.ID().String(),
		OrganizationID: u.OrganizationID().String(),
		RoleID:         u.RoleID().String(),
		Email:          u.Email().String(),
		FirstName:      u.FirstName(),
		LastName:       u.LastName(),
		CreatedAt:      u.CreatedAt(),
	}, nil
}

type memoryOrganizations struct{ m *memoryStore }

func (r memoryOrganizations) Create(_ context.Context, o *organization.Organization) error {
	r.m.orgs[o.ID()] = o
	return nil
}

func (r memoryOrganizations) FindByID(_ context.Context, id organization.ID) (*organization.Organization, error) {
	o, ok := r.m.orgs[id]
	if !ok {
		return nil, organization.ErrOrganizationNotFound
	}
	return o, nil
}

type memoryRoles struct{ m *memoryStore }

func (r memoryRoles) Create(_ context.Context, ro *role.Role) error {
	r.m.roles[ro.ID()] = ro
	return nil
}

func (r memoryRoles) FindByID(_ context.Context, id role.ID) (*role.Role, error) {
	ro, ok := r.m.roles[id]
	if !ok {
		return nil, role.ErrRoleNotFound
	}
	return ro, nil
}

type memoryUsers struct{ m *memoryStore }

func (r memoryUsers) Create(_ context.Context, u *user.User) error {
	if r.m.createErr != nil {
		return r.m.createErr
	}
	for _, existing := range r.m.users {
		if existing.Email() == u.Email() {
			return user.ErrEmailAlreadyExists
		}
	}
	r.m.users[u.ID()] = u
	return nil
}

func (r memoryUsers) FindByID(_ context.Context, id user.ID) (*user.User, error) {
	u, ok := r.m.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (r memoryUsers) FindByEmail(_ context.Context, email user.Email) (*user.User, error) {
	for _, u := range r.m.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, user.ErrUserNotFound
}

func newTestClient(t *testing.T, store *memoryStore) userv1.UserServiceClient {
	t.Helper()

	app := application.Application{
		Commands: application.Commands{CreateOrganization: command.NewCreateOrganizationHandler(store)},
		Queries:  application.Queries{GetUser: query.NewGetUserHandler(store)},
	}

	lis := bufconn.Listen(1024 * 1024)
	srv := grpclib.NewServer()
	userv1.RegisterUserServiceServer(srv, NewServer(app))
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	conn, err := grpclib.NewClient("passthrough:///bufnet",
		grpclib.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpclib.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	return userv1.NewUserServiceClient(conn)
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if got := status.Code(err); got != want {
		t.Fatalf("want code %s, got %s (%v)", want, got, err)
	}
}

func TestUserServiceServer(t *testing.T) {
	ctx := context.Background()
	store := newMemoryStore()
	client := newTestClient(t, store)

	created, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
		OrganizationName: "Acme", OwnerEmail: "Jane@Example.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	if err != nil {
		t.Fatalf("CreateOrganization: %v", err)
	}

	t.Run("GetUser returns the owner in the new organization", func(t *testing.T) {
		resp, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: created.GetOwnerUserId()})
		if err != nil {
			t.Fatal(err)
		}
		u := resp.GetUser()
		if u.GetId() != created.GetOwnerUserId() || u.GetEmail() != "jane@example.com" ||
			u.GetFirstName() != "Jane" || u.GetLastName() != "Doe" || u.GetCreatedAt().AsTime().IsZero() {
			t.Fatalf("got %+v", u)
		}
		if u.GetOrganizationId() != created.GetOrganizationId() {
			t.Errorf("organization_id = %q, want %q", u.GetOrganizationId(), created.GetOrganizationId())
		}

		roleID, err := role.ParseID(u.GetRoleId())
		if err != nil {
			t.Fatal(err)
		}
		if r, ok := store.roles[roleID]; !ok || !r.IsOwner() {
			t.Errorf("role_id %q is not the owner role", u.GetRoleId())
		}
	})

	t.Run("CreateOrganization with taken email is AlreadyExists", func(t *testing.T) {
		_, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
			OrganizationName: "Globex", OwnerEmail: "jane@example.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
		})
		assertCode(t, err, codes.AlreadyExists)
	})

	t.Run("CreateOrganization with invalid email is InvalidArgument", func(t *testing.T) {
		_, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
			OrganizationName: "Globex", OwnerEmail: "nope", OwnerFirstName: "Jane", OwnerLastName: "Doe",
		})
		assertCode(t, err, codes.InvalidArgument)
	})

	t.Run("CreateOrganization with empty name is InvalidArgument", func(t *testing.T) {
		_, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
			OrganizationName: " ", OwnerEmail: "john@example.com", OwnerFirstName: "John", OwnerLastName: "Doe",
		})
		assertCode(t, err, codes.InvalidArgument)
	})

	t.Run("GetUser with malformed ID is InvalidArgument", func(t *testing.T) {
		_, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: "not-a-uuid"})
		assertCode(t, err, codes.InvalidArgument)
	})

	t.Run("GetUser with unknown ID is NotFound", func(t *testing.T) {
		_, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: user.NewID().String()})
		assertCode(t, err, codes.NotFound)
	})
}

func TestUnexpectedErrorsAreHidden(t *testing.T) {
	store := newMemoryStore()
	store.createErr = errors.New("connection refused to db-host:5432")
	client := newTestClient(t, store)

	_, err := client.CreateOrganization(context.Background(), &userv1.CreateOrganizationRequest{
		OrganizationName: "Acme", OwnerEmail: "jane@example.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	assertCode(t, err, codes.Internal)
	if msg := status.Convert(err).Message(); msg != "internal error" {
		t.Fatalf("internal details leaked to caller: %q", msg)
	}
}
