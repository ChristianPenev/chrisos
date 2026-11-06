package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	graphqladapter "github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql/generated"
	"github.com/ChristianPenev/chrisos/pkg/httpserver"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gql := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graphqladapter.Resolver{},
	}))
	gql.AddTransport(transport.POST{})
	gql.Use(extension.Introspection{})

	mux := http.NewServeMux()
	mux.Handle("POST /graphql", gql)
	mux.Handle("GET /playground", playground.Handler("people-api", "/graphql"))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return httpserver.Run(ctx, ":8080", mux)
}
