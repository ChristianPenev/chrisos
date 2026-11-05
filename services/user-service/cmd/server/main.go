package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChristianPenev/chrisos/services/user/internal/database"
)

// shutdownTimeout is how long in-flight requests get to finish after SIGTERM.
// Must stay below the pod's terminationGracePeriodSeconds (k8s default: 30s).
const shutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// ctx is cancelled on SIGINT (Ctrl+C) or SIGTERM (sent by Kubernetes on pod stop).
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("PostgreSQL ping failed: %w", err)
	}

	log.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ChrisOS User Service OK")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Println("User Service listening on :8080")
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		// Server failed on its own (e.g. port already in use).
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		log.Println("Shutdown signal received, draining connections...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	// After Shutdown, ListenAndServe returns http.ErrServerClosed; anything else is a real error.
	if err := <-serverErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server: %w", err)
	}

	log.Println("User Service stopped")
	return nil
}
