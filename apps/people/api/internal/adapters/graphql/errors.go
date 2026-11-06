package graphql

import (
	"context"
	"errors"
	"log"
	"strings"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
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
	case errors.Is(err, app.ErrUserNotFound):
		setCode(gqlErr, codeNotFound)
	case errors.Is(err, app.ErrInvalidInput):
		setCode(gqlErr, codeBadUserInput)
	case errors.Is(err, app.ErrEmailAlreadyExists):
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
