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
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type memoryStore struct {
	users     map[user.ID]*user.User
	createErr error
}

func (m *memoryStore) Create(_ context.Context, u *user.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	for _, existing := range m.users {
		if existing.Email() == u.Email() {
			return user.ErrEmailAlreadyExists
		}
	}
	m.users[u.ID()] = u
	return nil
}

func (m *memoryStore) FindByID(_ context.Context, id user.ID) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (m *memoryStore) FindByEmail(_ context.Context, email user.Email) (*user.User, error) {
	for _, u := range m.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, user.ErrUserNotFound
}

func (m *memoryStore) UserByID(ctx context.Context, id user.ID) (query.UserView, error) {
	u, err := m.FindByID(ctx, id)
	if err != nil {
		return query.UserView{}, err
	}
	return query.UserView{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		FirstName: u.FirstName(),
		LastName:  u.LastName(),
		CreatedAt: u.CreatedAt(),
	}, nil
}

func newTestClient(t *testing.T, store *memoryStore) userv1.UserServiceClient {
	t.Helper()

	app := application.Application{
		Commands: application.Commands{RegisterUser: command.NewRegisterUserHandler(store)},
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
	client := newTestClient(t, &memoryStore{users: map[user.ID]*user.User{}})

	reg, err := client.RegisterUser(ctx, &userv1.RegisterUserRequest{
		Email: "Jane@Example.com", FirstName: "Jane", LastName: "Doe",
	})
	if err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}

	t.Run("GetUser returns the registered user", func(t *testing.T) {
		resp, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: reg.GetId()})
		if err != nil {
			t.Fatal(err)
		}
		u := resp.GetUser()
		if u.GetId() != reg.GetId() || u.GetEmail() != "jane@example.com" ||
			u.GetFirstName() != "Jane" || u.GetLastName() != "Doe" || u.GetCreatedAt().AsTime().IsZero() {
			t.Fatalf("got %+v", u)
		}
	})

	t.Run("RegisterUser with taken email is AlreadyExists", func(t *testing.T) {
		_, err := client.RegisterUser(ctx, &userv1.RegisterUserRequest{
			Email: "jane@example.com", FirstName: "Jane", LastName: "Doe",
		})
		assertCode(t, err, codes.AlreadyExists)
	})

	t.Run("RegisterUser with invalid email is InvalidArgument", func(t *testing.T) {
		_, err := client.RegisterUser(ctx, &userv1.RegisterUserRequest{
			Email: "nope", FirstName: "Jane", LastName: "Doe",
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
	client := newTestClient(t, &memoryStore{
		users:     map[user.ID]*user.User{},
		createErr: errors.New("connection refused to db-host:5432"),
	})

	_, err := client.RegisterUser(context.Background(), &userv1.RegisterUserRequest{
		Email: "jane@example.com", FirstName: "Jane", LastName: "Doe",
	})
	assertCode(t, err, codes.Internal)
	if msg := status.Convert(err).Message(); msg != "internal error" {
		t.Fatalf("internal details leaked to caller: %q", msg)
	}
}
