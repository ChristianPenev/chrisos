package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ChristianPenev/chrisos/services/user/internal/database"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("PostgreSQL ping failed: ", err)
	}

	log.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ChrisOS User Service OK")
	})

	log.Println("User Service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
