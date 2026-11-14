package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"testing"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
	"github.com/ChristianPenev/chrisos/pkg/batch"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/memory"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func newTestClient(t *testing.T, store *memory.Store) userv1.UserServiceClient {
	t.Helper()

	app := application.Application{
		Commands: application.Commands{CreateOrganization: command.NewCreateOrganizationHandler(store)},
		Queries: application.Queries{
			GetUser:   query.NewGetUserHandler(store),
			GetUsers:  query.NewGetUsersHandler(store),
			ListUsers: query.NewListUsersHandler(store),
			GetRoles:  query.NewGetRolesHandler(store),
		},
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
	store := memory.NewStore()
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
		if r, ok := store.Role(roleID); !ok || !r.IsOwner() {
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

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	client := newTestClient(t, store)

	acme, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
		OrganizationName: "Acme", OwnerEmail: "owner@acme.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
		OrganizationName: "Other", OwnerEmail: "owner@other.com", OwnerFirstName: "John", OwnerLastName: "Roe",
	}); err != nil {
		t.Fatal(err)
	}

	owner, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: acme.GetOwnerUserId()})
	if err != nil {
		t.Fatal(err)
	}
	orgID, _ := organization.ParseID(acme.GetOrganizationId())
	roleID, _ := role.ParseID(owner.GetUser().GetRoleId())

	want := []string{acme.GetOwnerUserId()}
	for i := range 4 {
		email, _ := user.NewEmail(fmt.Sprintf("member%d@acme.com", i))
		u, err := user.NewUser(orgID, roleID, email, "Member", "Acme")
		if err != nil {
			t.Fatal(err)
		}
		if err := store.Repositories().Users.Create(ctx, u); err != nil {
			t.Fatal(err)
		}
		want = append(want, u.ID().String())
	}
	slices.Sort(want)

	t.Run("pages through the organization's users", func(t *testing.T) {
		var got []string
		token := ""
		for range 10 {
			resp, err := client.ListUsers(ctx, &userv1.ListUsersRequest{
				OrganizationId: acme.GetOrganizationId(), PageSize: 2, PageToken: token,
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, u := range resp.GetUsers() {
				if u.GetOrganizationId() != acme.GetOrganizationId() || u.GetCreatedAt().AsTime().IsZero() {
					t.Fatalf("unexpected user %+v", u)
				}
				got = append(got, u.GetId())
			}
			token = resp.GetNextPageToken()
			if token == "" {
				break
			}
		}
		if !slices.Equal(got, want) {
			t.Fatalf("got  %v\nwant %v", got, want)
		}
	})

	invalid := []struct {
		name string
		req  *userv1.ListUsersRequest
	}{
		{"malformed organization ID", &userv1.ListUsersRequest{OrganizationId: "nope"}},
		{"negative page size", &userv1.ListUsersRequest{OrganizationId: acme.GetOrganizationId(), PageSize: -1}},
		{"garbage page token", &userv1.ListUsersRequest{OrganizationId: acme.GetOrganizationId(), PageToken: "garbage"}},
	}
	for _, tt := range invalid {
		t.Run(tt.name+" is InvalidArgument", func(t *testing.T) {
			_, err := client.ListUsers(ctx, tt.req)
			assertCode(t, err, codes.InvalidArgument)
		})
	}
}

func TestBatchGetRoles(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	client := newTestClient(t, store)

	created, err := client.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
		OrganizationName: "Acme", OwnerEmail: "owner@acme.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	owner, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: created.GetOwnerUserId()})
	if err != nil {
		t.Fatal(err)
	}
	ownerRoleID := owner.GetUser().GetRoleId()

	t.Run("returns known roles and leaves out unknown IDs", func(t *testing.T) {
		resp, err := client.BatchGetRoles(ctx, &userv1.BatchGetRolesRequest{
			Ids: []string{ownerRoleID, role.NewID().String(), ownerRoleID},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.GetRoles()) != 1 {
			t.Fatalf("got %d roles, want 1", len(resp.GetRoles()))
		}
		r := resp.GetRoles()[0]
		if r.GetId() != ownerRoleID || r.GetName() != role.OwnerRoleName || !r.GetOwner() ||
			r.GetOrganizationId() != created.GetOrganizationId() || r.GetCreatedAt().AsTime().IsZero() {
			t.Fatalf("got %+v", r)
		}
		if len(r.GetPermissions()) != len(role.AllPermissions()) {
			t.Errorf("permissions = %v, want all", r.GetPermissions())
		}
	})

	t.Run("empty request returns no roles", func(t *testing.T) {
		resp, err := client.BatchGetRoles(ctx, &userv1.BatchGetRolesRequest{})
		if err != nil || len(resp.GetRoles()) != 0 {
			t.Fatalf("got %v, %v", resp, err)
		}
	})

	t.Run("malformed ID is InvalidArgument", func(t *testing.T) {
		_, err := client.BatchGetRoles(ctx, &userv1.BatchGetRolesRequest{Ids: []string{"nope"}})
		assertCode(t, err, codes.InvalidArgument)
	})

	t.Run("too many IDs is InvalidArgument", func(t *testing.T) {
		ids := make([]string, batch.MaxIDs+1)
		for i := range ids {
			ids[i] = role.NewID().String()
		}
		_, err := client.BatchGetRoles(ctx, &userv1.BatchGetRolesRequest{Ids: ids})
		assertCode(t, err, codes.InvalidArgument)
	})
}

func TestUnexpectedErrorsAreHidden(t *testing.T) {
	store := memory.NewStore()
	store.UserCreateErr = errors.New("connection refused to db-host:5432")
	client := newTestClient(t, store)

	_, err := client.CreateOrganization(context.Background(), &userv1.CreateOrganizationRequest{
		OrganizationName: "Acme", OwnerEmail: "jane@example.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	assertCode(t, err, codes.Internal)
	if msg := status.Convert(err).Message(); msg != "internal error" {
		t.Fatalf("internal details leaked to caller: %q", msg)
	}
}
