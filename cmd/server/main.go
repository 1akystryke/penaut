package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"messenger/internal/handler"
	"messenger/internal/middleware"
	"messenger/internal/repository/postgres"
	"messenger/internal/service"
	"messenger/internal/storage/minio"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dbHostname := getenv("DOCKER_DB_NAME", "postgres")
	dbPort := getenv("DOCKER_DB_PORT", "5432")
	dbName := getenv("DB_NAME", "peanut")
	dbUser := getenv("DB_USER", "postgres")
	dbPass := getenv("DB_PASSWORD", "postgres")

	databaseURL := "postgres://" + dbUser + ":" + dbPass + "@" + dbHostname + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	fmt.Println("database url: " + databaseURL)
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	minioHostname := getenv("DOCKER_MINIO_NAME", "minio")
	minioPort := getenv("DOCKER_MINIO_PORT_WEB", "9000")
	minioAccessKey := getenv("MINIO_ACCESS_KEY", "minio")
	minioSecretKey := getenv("MINIO_SECRET_KEY", "minio")
	minioBucket := getenv("MINIO_BUCKET", "bucket")

	storage := minio.New(
		minioHostname+":"+minioPort,
		minioAccessKey,
		minioSecretKey,
		minioBucket,
		false,
	)

	store := postgres.NewStore(db)
	if getenv("INIT_DB", "true") == "true" {
		if err := store.Init(ctx); err != nil {
			log.Fatal(err)
		}
	}

	svc := service.New(store, store, store, store, store, store, *storage)

	h := handler.New(svc)
	handler := setupMiddleware(h, svc)

	addr := ":" + getenv("DOCKER_BE_PORT", "8080")
	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func setupMiddleware(h *handler.Handler, s *service.Service) http.Handler {
	authFunc := middleware.MakeAuthMiddleware(s, 30)
	var handler http.Handler = h.Routes()
	handler = middleware.Logging(handler)
	//handler = middleware.Auth(s, handler)
	handler = authFunc(handler)
	handler = middleware.CORS(handler)
	return handler
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
