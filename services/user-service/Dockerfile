FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /user-service ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /user-service ./user-service

EXPOSE 8080

USER 10001:10001

CMD ["./user-service"]