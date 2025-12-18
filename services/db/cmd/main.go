package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/postgres"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/redis"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/transport/grpc"
)

func main() {
	ctx := context.Background()

	postgresController := postgres.NewPostgresController()

	ttlSeconds, _ := strconv.Atoi(os.Getenv("REDIS_TTL_SECONDS"))
	redisCacheController, err := redis.NewRedisController(ctx, "redis:6379", ttlSeconds)
	if err != nil {
		log.Fatalf("failed to create redis cache controller: %v", err)
	}

	cacheDBRepository := app.NewCachedRepository(ctx, postgresController, redisCacheController)

	service := app.NewService(cacheDBRepository)

	server := grpc.NewServer(service)

	dbGrpcServerAddress := ":" + os.Getenv("DB_SERVICE_INTERNAL_PORT")
	if err := server.StartServer(dbGrpcServerAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
