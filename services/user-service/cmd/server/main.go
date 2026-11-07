package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	userv1 "github.com/ChristianPenev/chrisos/gen/go/user/v1"
	"github.com/ChristianPenev/chrisos/pkg/grpcserver"
	"github.com/ChristianPenev/chrisos/pkg/postgres"
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
			ListUsers: query.NewListUsersHandler(users),
			GetRoles:  query.NewGetRolesHandler(pgadapter.NewRoleReader(db)),
		},
	}

	return grpcserver.Run(ctx, ":9090", func(s grpc.ServiceRegistrar) {
		userv1.RegisterUserServiceServer(s, grpcadapter.NewServer(app))
	})
}
