package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
)

const maxQueryComplexity = 100

func NewHandler(application *app.Application) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{app: application},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(maxQueryComplexity))
	srv.SetErrorPresenter(presentError)
	return srv
}
