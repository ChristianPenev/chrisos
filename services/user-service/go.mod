module github.com/ChristianPenev/chrisos/services/user-service

go 1.27.3

require (
	github.com/99designs/gqlgen v0.17.98
	github.com/ChristianPenev/chrisos/gen/go v0.0.0
	github.com/ChristianPenev/chrisos/pkg v0.0.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.11.0
	github.com/vektah/gqlparser/v2 v2.5.64
	golang.org/x/sync v0.24.0
	google.golang.org/grpc v1.85.0
	google.golang.org/protobuf v1.36.12
	gorm.io/driver/postgres v1.6.3
	gorm.io/gorm v1.31.2
)

require (
	github.com/agnivade/levenshtein v1.2.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coder/websocket v1.8.15 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose/v3 v3.29.0 // indirect
	github.com/sethvargo/go-retry v0.5.0 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	github.com/vikstrous/dataloadgen v0.0.10 // indirect
	go.opentelemetry.io/otel v1.47.0 // indirect
	go.opentelemetry.io/otel/trace v1.47.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
)

replace github.com/ChristianPenev/chrisos/pkg => ../../pkg

replace github.com/ChristianPenev/chrisos/gen/go => ../../gen/go
