package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/postgres"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/redis"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/transport/grpc"
)

func main() {
	logPath := os.Getenv("LOG_FILE_PATH")

	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal("open log file error:", err)
	}
	defer func() { _ = f.Close() }()

	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx := context.Background()

	postgresController := postgres.NewPostgresController()

	ttlSeconds, _ := strconv.Atoi(os.Getenv("REDIS_TTL_SECONDS"))
	redisCacheController, err := redis.NewRedisController(ctx, "redis:6379", ttlSeconds)
	if err != nil {
		log.Fatalf("failed to create redis cache controller: %v\n", err)
	}

	cacheDBRepository := app.NewCachedRepository(postgresController, redisCacheController)

	service := app.NewService(cacheDBRepository)

	server := grpc.NewServer(service)

	dbGrpcServerAddress := ":" + os.Getenv("DB_SERVICE_INTERNAL_PORT")
	if err := server.StartServer(dbGrpcServerAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
