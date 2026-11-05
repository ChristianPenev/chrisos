package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ChristianPenev/chrisos/pkg/httpserver"
	"github.com/ChristianPenev/chrisos/pkg/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// ctx is cancelled on SIGINT (Ctrl+C) or SIGTERM (sent by Kubernetes on pod stop).
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ChrisOS User Service OK")
	})

	return httpserver.Run(ctx, ":8080", mux)
}
