package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func toStatus(err error) error {
	switch {
	case errors.Is(err, user.ErrInvalidID),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidFirstName),
		errors.Is(err, user.ErrInvalidLastName),
		errors.Is(err, organization.ErrInvalidID),
		errors.Is(err, organization.ErrInvalidName),
		errors.Is(err, role.ErrInvalidID),
		errors.Is(err, role.ErrInvalidName),
		errors.Is(err, role.ErrUnknownPermission):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, user.ErrUserNotFound),
		errors.Is(err, organization.ErrOrganizationNotFound),
		errors.Is(err, role.ErrRoleNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, user.ErrEmailAlreadyExists),
		errors.Is(err, role.ErrRoleNameTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, role.ErrOwnerRoleLocked):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()
	default:
		log.Printf("internal error: %v", err)
		return status.Error(codes.Internal, "internal error")
	}
}
