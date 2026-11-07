package userservice

import (
	"context"
	"errors"
	"math"
	"net"
	"slices"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
)

type fakeUserService struct {
	userv1.UnimplementedUserServiceServer
	err   error
	delay time.Duration
	got   *userv1.CreateOrganizationRequest
	list  *userv1.ListUsersRequest
}

func (f *fakeUserService) BatchGetRoles(ctx context.Context, req *userv1.BatchGetRolesRequest) (*userv1.BatchGetRolesResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	var roles []*userv1.Role
	for _, id := range req.GetIds() {
		roles = append(roles, &userv1.Role{
			Id: id, OrganizationId: "org-1", Name: "Support", Owner: false,
			Permissions: []string{"members.invite"},
			CreatedAt:   timestamppb.New(time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)),
		})
	}
	return &userv1.BatchGetRolesResponse{Roles: roles}, nil
}

func (f *fakeUserService) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	f.list = req
	if f.err != nil {
		return nil, f.err
	}
	return &userv1.ListUsersResponse{
		Users: []*userv1.User{
			{Id: "user-1", OrganizationId: req.GetOrganizationId(), Email: "a@example.com"},
			{Id: "user-2", OrganizationId: req.GetOrganizationId(), Email: "b@example.com"},
		},
		NextPageToken: "next",
	}, nil
}

func (f *fakeUserService) CreateOrganization(ctx context.Context, req *userv1.CreateOrganizationRequest) (*userv1.CreateOrganizationResponse, error) {
	f.got = req
	if f.err != nil {
		return nil, f.err
	}
	return &userv1.CreateOrganizationResponse{OrganizationId: "org-1", OwnerUserId: "user-1"}, nil
}

func (f *fakeUserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	select {
	case <-time.After(f.delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if f.err != nil {
		return nil, f.err
	}
	return &userv1.GetUserResponse{User: &userv1.User{
		Id:             req.GetId(),
		OrganizationId: "org-1",
		RoleId:         "role-1",
		Email:          "jane@example.com",
		FirstName:      "Jane",
		LastName:       "Doe",
		CreatedAt:      timestamppb.New(time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)),
	}}, nil
}

func newTestClient(t *testing.T, fake *fakeUserService) *Client {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	userv1.RegisterUserServiceServer(srv, fake)
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	return NewClient(conn)
}

func TestCreateOrganization(t *testing.T) {
	fake := &fakeUserService{}
	got, err := newTestClient(t, fake).CreateOrganization(context.Background(), app.CreateOrganizationInput{
		Name: "Acme", OwnerEmail: "jane@example.com", OwnerFirstName: "Jane", OwnerLastName: "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != (app.CreatedOrganization{OrganizationID: "org-1", OwnerUserID: "user-1"}) {
		t.Errorf("got %+v", got)
	}
	if fake.got.GetOrganizationName() != "Acme" || fake.got.GetOwnerEmail() != "jane@example.com" ||
		fake.got.GetOwnerFirstName() != "Jane" || fake.got.GetOwnerLastName() != "Doe" {
		t.Errorf("request not mapped: %+v", fake.got)
	}
}

func TestGetUser(t *testing.T) {
	u, err := newTestClient(t, &fakeUserService{}).GetUser(context.Background(), "user-1")
	if err != nil {
		t.Fatal(err)
	}
	want := app.User{
		ID: "user-1", OrganizationID: "org-1", RoleID: "role-1",
		Email: "jane@example.com", FirstName: "Jane", LastName: "Doe",
		CreatedAt: time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC),
	}
	if *u != want {
		t.Errorf("got %+v, want %+v", *u, want)
	}
}

func TestListUsers(t *testing.T) {
	fake := &fakeUserService{}
	page, err := newTestClient(t, fake).ListUsers(context.Background(), app.ListUsersInput{
		OrganizationID: "org-1", First: 2, After: "abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	if fake.list.GetOrganizationId() != "org-1" || fake.list.GetPageSize() != 2 || fake.list.GetPageToken() != "abc" {
		t.Errorf("request not mapped: %+v", fake.list)
	}
	if len(page.Users) != 2 || page.Users[1].Email != "b@example.com" || page.Users[0].OrganizationID != "org-1" {
		t.Errorf("users not mapped: %+v", page.Users)
	}
	if page.NextCursor != "next" {
		t.Errorf("next cursor = %q", page.NextCursor)
	}
}

func TestListUsersClampsHugePageSize(t *testing.T) {
	fake := &fakeUserService{}
	if _, err := newTestClient(t, fake).ListUsers(context.Background(), app.ListUsersInput{
		OrganizationID: "org-1", First: math.MaxInt32 + 1,
	}); err != nil {
		t.Fatal(err)
	}
	if fake.list.GetPageSize() != math.MaxInt32 {
		t.Errorf("page size = %d, want %d", fake.list.GetPageSize(), int32(math.MaxInt32))
	}
}

func TestGetRoles(t *testing.T) {
	roles, err := newTestClient(t, &fakeUserService{}).GetRoles(context.Background(), []string{"role-1", "role-2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 2 {
		t.Fatalf("got %d roles", len(roles))
	}
	r := roles[1]
	if r.ID != "role-2" || r.OrganizationID != "org-1" || r.Name != "Support" || r.IsOwner ||
		!slices.Equal(r.Permissions, []string{"members.invite"}) ||
		!r.CreatedAt.Equal(time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)) {
		t.Errorf("role not mapped: %+v", r)
	}
}

func TestErrorTranslation(t *testing.T) {
	tests := []struct {
		name    string
		grpcErr error
		want    error
	}{
		{"NotFound", status.Error(codes.NotFound, "user not found"), app.ErrUserNotFound},
		{"InvalidArgument", status.Error(codes.InvalidArgument, `invalid email: "nope"`), app.ErrInvalidInput},
		{"AlreadyExists", status.Error(codes.AlreadyExists, "email already exists"), app.ErrEmailAlreadyExists},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(t, &fakeUserService{err: tt.grpcErr}).GetUser(context.Background(), "user-1")
			if !errors.Is(err, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, err)
			}
		})
	}

	t.Run("InvalidArgument keeps the reason", func(t *testing.T) {
		_, err := newTestClient(t, &fakeUserService{
			err: status.Error(codes.InvalidArgument, `invalid email: "nope"`),
		}).GetUser(context.Background(), "user-1")
		if err.Error() != `invalid input: invalid email: "nope"` {
			t.Fatalf("got %q", err)
		}
	})

	t.Run("unexpected codes are not mapped to app errors", func(t *testing.T) {
		_, err := newTestClient(t, &fakeUserService{err: status.Error(codes.Internal, "boom")}).GetUser(context.Background(), "user-1")
		if err == nil || errors.Is(err, app.ErrUserNotFound) || errors.Is(err, app.ErrInvalidInput) {
			t.Fatalf("got %v", err)
		}
	})
}

func TestDeadline(t *testing.T) {
	client := newTestClient(t, &fakeUserService{delay: time.Second})
	client.timeout = 50 * time.Millisecond

	start := time.Now()
	_, err := client.GetUser(context.Background(), "user-1")
	if status.Code(errors.Unwrap(err)) != codes.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("call took %s, deadline not applied", elapsed)
	}
}
