module github.com/ChristianPenev/chrisos/apps/people/api

go 1.27.3

require (
	github.com/99designs/gqlgen v0.17.98
	github.com/ChristianPenev/chrisos/gen/go v0.0.0
	github.com/ChristianPenev/chrisos/pkg v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.64
	google.golang.org/grpc v1.85.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/agnivade/levenshtein v1.2.2 // indirect
	github.com/coder/websocket v1.8.15 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
)

replace github.com/ChristianPenev/chrisos/pkg => ../../../pkg

replace github.com/ChristianPenev/chrisos/gen/go => ../../../gen/go
