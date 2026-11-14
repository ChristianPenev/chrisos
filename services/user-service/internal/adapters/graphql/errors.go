package graphql

import (
	"context"
	"errors"
	"log"
	"strings"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/ChristianPenev/chrisos/pkg/batch"
	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

const (
	codeNotFound     = "NOT_FOUND"
	codeBadUserInput = "BAD_USER_INPUT"
	codeConflict     = "CONFLICT"
	codeInternal     = "INTERNAL"
)

func presentError(ctx context.Context, err error) *gqlerror.Error {
	gqlErr := gql.DefaultErrorPresenter(ctx, err)

	switch {
	case errors.Is(err, user.ErrInvalidID),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidFirstName),
		errors.Is(err, user.ErrInvalidLastName),
		errors.Is(err, organization.ErrInvalidID),
		errors.Is(err, organization.ErrInvalidName),
		errors.Is(err, role.ErrInvalidID),
		errors.Is(err, role.ErrInvalidName),
		errors.Is(err, role.ErrUnknownPermission),
		errors.Is(err, pagination.ErrInvalidPageSize),
		errors.Is(err, pagination.ErrInvalidCursor),
		errors.Is(err, batch.ErrTooManyIDs):
		setCode(gqlErr, codeBadUserInput)
	case errors.Is(err, user.ErrUserNotFound),
		errors.Is(err, organization.ErrOrganizationNotFound),
		errors.Is(err, role.ErrRoleNotFound):
		setCode(gqlErr, codeNotFound)
	case errors.Is(err, user.ErrEmailAlreadyExists),
		errors.Is(err, role.ErrRoleNameTaken):
		setCode(gqlErr, codeConflict)
	case gqlErr.Extensions["code"] != nil, isIntrospectionField(ctx):
	default:
		log.Printf("graphql %s: %v", gqlErr.Path, err)
		gqlErr.Message = "internal error"
		setCode(gqlErr, codeInternal)
	}
	return gqlErr
}

func isIntrospectionField(ctx context.Context) bool {
	fc := gql.GetFieldContext(ctx)
	return fc != nil && strings.HasPrefix(fc.Field.Name, "__")
}

func setCode(e *gqlerror.Error, code string) {
	if e.Extensions == nil {
		e.Extensions = map[string]any{}
	}
	e.Extensions["code"] = code
}
