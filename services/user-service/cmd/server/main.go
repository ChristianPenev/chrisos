package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/playground"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
	"github.com/ChristianPenev/chrisos/pkg/grpcserver"
	"github.com/ChristianPenev/chrisos/pkg/httpserver"
	"github.com/ChristianPenev/chrisos/pkg/postgres"
	graphqladapter "github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/graphql"
	grpcadapter "github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/grpc"
	pgadapter "github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/postgres"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/postgres/migrations"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbCfg, err := postgres.ConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := postgres.Connect(dbCfg)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	log.Println("Connected to PostgreSQL")

	if err := postgres.Migrate(ctx, sqlDB, migrations.FS); err != nil {
		return err
	}

	users := pgadapter.NewUserReader(db)

	app := application.Application{
		Commands: application.Commands{
			CreateOrganization: command.NewCreateOrganizationHandler(pgadapter.NewTransactor(db)),
		},
		Queries: application.Queries{
			GetUser:   query.NewGetUserHandler(users),
			GetUsers:  query.NewGetUsersHandler(users),
			ListUsers: query.NewListUsersHandler(users),
			GetRoles:  query.NewGetRolesHandler(pgadapter.NewRoleReader(db)),
		},
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return grpcserver.Run(ctx, ":9090", func(s grpc.ServiceRegistrar) {
			userv1.RegisterUserServiceServer(s, grpcadapter.NewServer(app))
		})
	})
	g.Go(func() error {
		return httpserver.Run(ctx, ":8080", newHTTPHandler(app))
	})
	return g.Wait()
}

func newHTTPHandler(app application.Application) http.Handler {
	devTools := os.Getenv("GRAPHQL_DEV_TOOLS") == "true"

	mux := http.NewServeMux()
	mux.Handle("POST /graphql", graphqladapter.NewHandler(app, devTools))
	if devTools {
		mux.Handle("GET /graphql", playground.ApolloSandboxHandler("user-service", "/graphql",
			playground.WithApolloSandboxEndpointIsEditable(false),
		))
		log.Println("GraphQL introspection and Apollo Sandbox enabled on GET /graphql")
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}
