# Build context: repo root (see Tiltfile), so the shared ./pkg and ./gen/go modules are visible.
FROM golang:1.27.3-alpine AS builder

# Each service builds standalone via the `replace` in its go.mod;
# go.work is for local dev only.
ENV GOWORK=off

WORKDIR /src

# Module files first so `go mod download` is cached until dependencies change.
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY gen/go/go.mod gen/go/go.sum ./gen/go/
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/

WORKDIR /src/services/user-service
RUN go mod download

COPY pkg/ /src/pkg/
COPY gen/go/ /src/gen/go/
COPY services/user-service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /user-service ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /user-service ./user-service

EXPOSE 9090

USER 10001:10001

CMD ["./user-service"]
