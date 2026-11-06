package userservice

import (
	"context"
	"fmt"
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

func (c *Client) RegisterUser(ctx context.Context, in app.RegisterUserInput) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.users.RegisterUser(ctx, &userv1.RegisterUserRequest{
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	})
	if err != nil {
		return "", translate(err)
	}
	return resp.GetId(), nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*app.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.users.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, translate(err)
	}

	u := resp.GetUser()
	return &app.User{
		ID:        u.GetId(),
		Email:     u.GetEmail(),
		FirstName: u.GetFirstName(),
		LastName:  u.GetLastName(),
		CreatedAt: u.GetCreatedAt().AsTime(),
	}, nil
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
