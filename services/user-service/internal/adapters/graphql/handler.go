package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/ChristianPenev/chrisos/pkg/dataloader"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
)

const maxQueryComplexity = 100

func NewHandler(app application.Application, introspection bool) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{app: app},
	}))
	srv.AddTransport(transport.POST{})
	if introspection {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.FixedComplexityLimit(maxQueryComplexity))
	srv.SetErrorPresenter(presentError)

	return dataloader.Middleware(func() *Loaders { return newLoaders(app) })(srv)
}
