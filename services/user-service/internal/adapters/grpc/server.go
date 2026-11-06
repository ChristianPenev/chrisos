package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
)

type Server struct {
	userv1.UnimplementedUserServiceServer
	app application.Application
}

func NewServer(app application.Application) *Server {
	return &Server{app: app}
}

func (s *Server) RegisterUser(ctx context.Context, req *userv1.RegisterUserRequest) (*userv1.RegisterUserResponse, error) {
	id, err := s.app.Commands.RegisterUser.Handle(ctx, command.RegisterUser{
		Email:     req.GetEmail(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
	})
	if err != nil {
		return nil, toStatus(err)
	}
	return &userv1.RegisterUserResponse{Id: id.String()}, nil
}

func (s *Server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	v, err := s.app.Queries.GetUser.Handle(ctx, query.GetUser{ID: req.GetId()})
	if err != nil {
		return nil, toStatus(err)
	}
	return &userv1.GetUserResponse{User: &userv1.User{
		Id:        v.ID,
		Email:     v.Email,
		FirstName: v.FirstName,
		LastName:  v.LastName,
		CreatedAt: timestamppb.New(v.CreatedAt),
	}}, nil
}
