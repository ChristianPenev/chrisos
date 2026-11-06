package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func toStatus(err error) error {
	switch {
	case errors.Is(err, user.ErrInvalidID),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidFirstName),
		errors.Is(err, user.ErrInvalidLastName):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, user.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()
	default:
		log.Printf("internal error: %v", err)
		return status.Error(codes.Internal, "internal error")
	}
}
