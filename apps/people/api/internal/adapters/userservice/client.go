package userservice

import (
	"context"
	"fmt"
	"math"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
)

const defaultTimeout = 3 * time.Second

type Client struct {
	users   userv1.UserServiceClient
	timeout time.Duration
}

var _ app.UserService = (*Client)(nil)

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{users: userv1.NewUserServiceClient(conn), timeout: defaultTimeout}
}

func (c *Client) CreateOrganization(ctx context.Context, in app.CreateOrganizationInput) (app.CreatedOrganization, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.users.CreateOrganization(ctx, &userv1.CreateOrganizationRequest{
		OrganizationName: in.Name,
		OwnerEmail:       in.OwnerEmail,
		OwnerFirstName:   in.OwnerFirstName,
		OwnerLastName:    in.OwnerLastName,
	})
	if err != nil {
		return app.CreatedOrganization{}, translate(err)
	}
	return app.CreatedOrganization{
		OrganizationID: resp.GetOrganizationId(),
		OwnerUserID:    resp.GetOwnerUserId(),
	}, nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*app.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.users.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, translate(err)
	}

	return fromProtoUser(resp.GetUser()), nil
}

func (c *Client) ListUsers(ctx context.Context, in app.ListUsersInput) (*app.UserPage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.users.ListUsers(ctx, &userv1.ListUsersRequest{
		OrganizationId: in.OrganizationID,
		PageSize:       int32(min(in.First, math.MaxInt32)),
		PageToken:      in.After,
	})
	if err != nil {
		return nil, translate(err)
	}

	users := make([]*app.User, len(resp.GetUsers()))
	for i, u := range resp.GetUsers() {
		users[i] = fromProtoUser(u)
	}
	return &app.UserPage{Users: users, NextCursor: resp.GetNextPageToken()}, nil
}

func (c *Client) GetRoles(ctx context.Context, ids []string) ([]*app.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.users.BatchGetRoles(ctx, &userv1.BatchGetRolesRequest{Ids: ids})
	if err != nil {
		return nil, translate(err)
	}

	roles := make([]*app.Role, len(resp.GetRoles()))
	for i, r := range resp.GetRoles() {
		roles[i] = &app.Role{
			ID:             r.GetId(),
			OrganizationID: r.GetOrganizationId(),
			Name:           r.GetName(),
			IsOwner:        r.GetOwner(),
			Permissions:    r.GetPermissions(),
			CreatedAt:      r.GetCreatedAt().AsTime(),
		}
	}
	return roles, nil
}

func fromProtoUser(u *userv1.User) *app.User {
	return &app.User{
		ID:             u.GetId(),
		OrganizationID: u.GetOrganizationId(),
		RoleID:         u.GetRoleId(),
		Email:          u.GetEmail(),
		FirstName:      u.GetFirstName(),
		LastName:       u.GetLastName(),
		CreatedAt:      u.GetCreatedAt().AsTime(),
	}
}

func translate(err error) error {
	st := status.Convert(err)
	switch st.Code() {
	case codes.NotFound:
		return app.ErrUserNotFound
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %s", app.ErrInvalidInput, st.Message())
	case codes.AlreadyExists:
		return app.ErrEmailAlreadyExists
	default:
		return fmt.Errorf("user-service: %w", err)
	}
}
