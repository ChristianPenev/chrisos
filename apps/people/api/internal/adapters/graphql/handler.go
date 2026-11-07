package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
	"github.com/ChristianPenev/chrisos/pkg/dataloader"
)

const maxQueryComplexity = 100

func NewHandler(application *app.Application, introspection bool) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{app: application},
	}))
	srv.AddTransport(transport.POST{})
	if introspection {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.FixedComplexityLimit(maxQueryComplexity))
	srv.SetErrorPresenter(presentError)

	return dataloader.Middleware(func() *Loaders { return newLoaders(application) })(srv)
}
