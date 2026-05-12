package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"messenger/internal/handler"
	"messenger/internal/middleware"
	"messenger/internal/repository/postgres"
	"messenger/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	databaseURL := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/messenger?sslmode=disable")
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := postgres.NewStore(db)
	if getenv("INIT_DB", "true") == "true" {
		if err := store.Init(ctx); err != nil {
			log.Fatal(err)
		}
	}

	svc := service.New(store, store, store, store)
	h := handler.New(svc)

	addr := getenv("ADDR", ":8080")
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.Logging(h.Routes())))
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
