package application

import (
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterUser command.RegisterUserHandler
}

type Queries struct {
	GetUser query.GetUserHandler
}
