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

func (s *Server) CreateOrganization(ctx context.Context, req *userv1.CreateOrganizationRequest) (*userv1.CreateOrganizationResponse, error) {
	res, err := s.app.Commands.CreateOrganization.Handle(ctx, command.CreateOrganization{
		Name:      req.GetOrganizationName(),
		Email:     req.GetOwnerEmail(),
		FirstName: req.GetOwnerFirstName(),
		LastName:  req.GetOwnerLastName(),
	})
	if err != nil {
		return nil, toStatus(err)
	}
	return &userv1.CreateOrganizationResponse{
		OrganizationId: res.OrganizationID.String(),
		OwnerUserId:    res.UserID.String(),
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	v, err := s.app.Queries.GetUser.Handle(ctx, query.GetUser{ID: req.GetId()})
	if err != nil {
		return nil, toStatus(err)
	}
	return &userv1.GetUserResponse{User: &userv1.User{
		Id:             v.ID,
		OrganizationId: v.OrganizationID,
		RoleId:         v.RoleID,
		Email:          v.Email,
		FirstName:      v.FirstName,
		LastName:       v.LastName,
		CreatedAt:      timestamppb.New(v.CreatedAt),
	}}, nil
}
