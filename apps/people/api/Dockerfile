FROM golang:1.27.3-alpine AS builder

ENV GOWORK=off

WORKDIR /src

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY gen/go/go.mod gen/go/go.sum ./gen/go/
COPY apps/people/api/go.mod apps/people/api/go.sum ./apps/people/api/

WORKDIR /src/apps/people/api
RUN go mod download

COPY pkg/ /src/pkg/
COPY gen/go/ /src/gen/go/
COPY apps/people/api/ ./

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /people-api ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /people-api ./people-api

EXPOSE 8080

USER 10001:10001

CMD ["./people-api"]
