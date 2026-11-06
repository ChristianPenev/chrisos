package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/playground"

	graphqladapter "github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/graphql"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/adapters/userservice"
	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
	"github.com/ChristianPenev/chrisos/pkg/grpcclient"
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

	userServiceAddr := envOr("USER_SERVICE_ADDR", "user-service:9090")
	conn, err := grpcclient.New(userServiceAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	application := app.New(userservice.NewClient(conn))
	devTools := os.Getenv("GRAPHQL_DEV_TOOLS") == "true"

	mux := http.NewServeMux()
	mux.Handle("POST /graphql", graphqladapter.NewHandler(application, devTools))
	if devTools {
		mux.Handle("GET /graphql", playground.ApolloSandboxHandler("people-api", "/graphql",
			playground.WithApolloSandboxEndpointIsEditable(false),
		))
		log.Println("GraphQL introspection and Apollo Sandbox enabled on GET /graphql")
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Calling user-service at %s", userServiceAddr)
	return httpserver.Run(ctx, ":8080", mux)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
